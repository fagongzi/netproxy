package proxy

import (
	"sync"

	"math/rand"

	"github.com/CodisLabs/codis/pkg/utils/log"
	"github.com/fagongzi/goetty"
	"github.com/fagongzi/netproxy/pkg/conf"
	"github.com/fagongzi/netproxy/pkg/util"
	"github.com/labstack/echo"

	"encoding/json"
	"io"
	"sort"
	"time"
)

// Ctl ctl lossRate and so on
type Ctl struct {
	Address string   `json:"address"`
	In      *CtlUnit `json:"in"`
	Out     *CtlUnit `json:"out"`
}

// CtlUnit CtlUnit
type CtlUnit struct {
	LossRate int `json:"lossRate"`
	DelayMs  int `json:"delayMs"`
}

// UnMarshalCtlFromReader UnMarshalCtlFromReader
func UnMarshalCtlFromReader(r io.Reader) (*Ctl, error) {
	v := &Ctl{}

	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)

	if nil != err {
		return nil, err
	}

	return v, nil
}

// Marshal marshal
func (c *Ctl) Marshal() []byte {
	d, _ := json.Marshal(c)
	return d
}

// Proxy proxy
type Proxy struct {
	sync.RWMutex
	cnf       *conf.Conf
	apiServer *echo.Echo
	wg        *sync.WaitGroup
	servers   []*TCPServer

	ctls map[string]*Ctl
}

// NewProxy factory method
func NewProxy(cnf *conf.Conf) *Proxy {
	return &Proxy{
		cnf:       cnf,
		apiServer: echo.New(),
		ctls:      make(map[string]*Ctl),
		wg:        &sync.WaitGroup{},
		servers:   make([]*TCPServer, len(cnf.Proxys)),
	}
}

// Start start server
func (p *Proxy) Start() {
	go p.startAPIServer()

	p.wg.Add(len(p.cnf.Proxys))
	for index, proxy := range p.cnf.Proxys {
		go func(proxy *conf.Proxy, index int) {
			server := &TCPServer{
				proxy:  proxy,
				server: goetty.NewServer(proxy.Src, DECODER, ENCODER, goetty.NewInt64IdGenerator()),
				p:      p,
			}
			p.servers[index] = server
			server.start()
			p.wg.Done()
		}(proxy, index)
	}

	p.wg.Wait()
}

// Stop stop server
func (p *Proxy) Stop() {
	for _, server := range p.servers {
		server.stop()
	}
}

// GetAllClients GetAllClients
func (p *Proxy) GetAllClients() []string {
	p.RLock()
	clients := make([]string, len(p.ctls))
	index := 0
	for key := range p.ctls {
		clients[index] = key
		index++
	}
	p.RUnlock()

	sort.Strings(clients)

	return clients
}

// UpdateCtl UpdateCtl
func (p *Proxy) UpdateCtl(ctl *Ctl) {
	p.Lock()
	p.ctls[ctl.Address] = ctl
	p.Unlock()
}

func (p *Proxy) addClientCtl(addr string) {
	p.Lock()
	p.ctls[addr] = &Ctl{
		Address: addr,
		In:      &CtlUnit{},
		Out:     &CtlUnit{},
	}
	p.Unlock()
}

func (p *Proxy) deleteClientCtl(addr string) {
	p.Lock()
	delete(p.ctls, addr)
	p.Unlock()
}

func (p *Proxy) getCtl(addr string) *Ctl {
	p.RLock()
	ctl := p.ctls[addr]
	p.RUnlock()

	return ctl
}

// TCPServer TCPServer
type TCPServer struct {
	proxy  *conf.Proxy
	p      *Proxy
	server *goetty.Server
}

func (t *TCPServer) start() {
	log.Infof("proxy <%s> to <%s>", t.proxy.Src, t.proxy.Target)
	t.server.Serve(t.doServe)
}

func (t *TCPServer) stop() {
	t.server.Stop()
}

func (t *TCPServer) doServe(session goetty.IOSession) error {
	t.p.addClientCtl(session.RemoteAddr())
	defer t.p.deleteClientCtl(session.RemoteAddr())

	var data interface{}
	var err error

	// client connected, make a connection to target
	conn := goetty.NewConnector(t.createGoettyConf(), DECODER, ENCODER)
	_, err = conn.Connect()

	if err != nil {
		log.InfoErrorf(err, "Connect to <%s> failure.", t.proxy.Target)
		return err
	}

	defer conn.Close()

	// read loop from target
	go func() {
		for {
			data, err := conn.Read()
			if err != nil {
				return
			}

			bytes, _ := data.([]byte)

			// write bytes to client
			ctl := t.p.getCtl(session.RemoteAddr())

			if 0 == ctl.In.LossRate {
				t.doWriteToClient(bytes, session, ctl.In)
			} else {
				if rand.Intn(100) > ctl.In.LossRate {
					t.doWriteToClient(bytes, session, ctl.In)
				} else {
					log.Infof("Loss write to <%s>", bytes, session.RemoteAddr())
				}
			}
		}
	}()

	for {
		data, err = session.Read()

		if err != nil {
			log.InfoErrorf(err, "Read from client<%s> failure.", session.RemoteAddr())
			break
		} else {
			bytes, _ := data.([]byte)

			// write to target
			ctl := t.p.getCtl(session.RemoteAddr())

			if 0 == ctl.Out.LossRate {
				t.doWrite(bytes, conn, ctl.Out)
			} else {
				if rand.Intn(100) > ctl.Out.LossRate {
					t.doWrite(bytes, conn, ctl.Out)
				} else {
					log.Infof("Loss write <%+v> to <%s>", bytes, t.proxy.Target)
				}
			}
		}
	}

	return err
}

func (t *TCPServer) writeTimeout(addr string, conn *goetty.Connector) {
	log.Warnf("Conn<%s> write timeout.", addr)
}

func (t *TCPServer) createGoettyConf() *goetty.Conf {
	return &goetty.Conf{
		Addr:                   t.proxy.Target,
		TimeWheel:              util.GetTimeWheel(),
		TimeoutWrite:           time.Second * time.Duration(t.proxy.TimeoutWrite),
		TimeoutConnectToServer: time.Second * time.Duration(t.proxy.TimeoutConnect),
		WriteTimeoutFn:         t.writeTimeout,
	}
}

func (t *TCPServer) doWrite(bytes []byte, conn *goetty.Connector, ctl *CtlUnit) {
	if ctl.DelayMs > 0 {
		log.Infof("Delay <%d>ms write to <%s>", ctl.DelayMs, t.proxy.Target)
		time.Sleep(time.Millisecond * time.Duration(ctl.DelayMs))
	}

	conn.Write(bytes)
	log.Infof("Write <%+v> to <%s>", bytes, t.proxy.Target)
}

func (t *TCPServer) doWriteToClient(bytes []byte, session goetty.IOSession, ctl *CtlUnit) {
	if ctl.DelayMs > 0 {
		log.Infof("Delay <%d>ms write to client<%s>", ctl.DelayMs, session.RemoteAddr())
		time.Sleep(time.Millisecond * time.Duration(ctl.DelayMs))
	}

	session.Write(bytes)
	log.Infof("Write <%+v> to client<%s>", bytes, session.RemoteAddr())
}
