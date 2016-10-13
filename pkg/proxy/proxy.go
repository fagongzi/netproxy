package proxy

import (
	"sync"

	"math/rand"

	"github.com/CodisLabs/codis/pkg/utils/log"
	"github.com/fagongzi/goetty"
	"github.com/fagongzi/netproxy/pkg/conf"
	"github.com/fagongzi/netproxy/pkg/util"
	"github.com/labstack/echo"

	"time"
)

// Proxy proxy
type Proxy struct {
	sync.RWMutex
	cnf       *conf.Conf
	apiServer *echo.Echo
	servers   map[string]*TCPServer
}

// NewProxy factory method
func NewProxy(cnf *conf.Conf) *Proxy {
	return &Proxy{
		cnf:       cnf,
		apiServer: echo.New(),
		servers:   make(map[string]*TCPServer),
	}
}

// Start start server
func (p *Proxy) Start() {
	for index, proxyUnit := range p.cnf.Units {
		go func(proxyUnit *conf.ProxyUnit, index int) {
			server := &TCPServer{
				proxyUnit: proxyUnit,
				p:         p,
			}
			p.servers[proxyUnit.Src] = server
			server.start()
		}(proxyUnit, index)
	}

	p.startAPIServer()
}

// Pause pause proxy listen
func (p *Proxy) Pause(addr string) {
	for _, server := range p.servers {
		if addr == server.proxyUnit.Src {
			server.pause()
		}
	}
}

// Resume resume proxy listen
func (p *Proxy) Resume(addr string) {
	for _, server := range p.servers {
		if addr == server.proxyUnit.Src {
			server.resume()
		}
	}
}

// Stop stop server
func (p *Proxy) Stop() {
	for _, server := range p.servers {
		server.stop()
	}
}

// UpdateCtl UpdateCtl
func (p *Proxy) UpdateCtl(ctl *conf.Ctl) {
	p.Lock()
	p.servers[ctl.Address].proxyUnit.Ctl.CopyFrom(ctl)
	p.Unlock()
}

// TCPServer TCPServer
type TCPServer struct {
	sync.RWMutex
	proxyUnit *conf.ProxyUnit
	p         *Proxy
	server    *goetty.Server
	paused    bool
}

func (t *TCPServer) start() {
	log.Infof("proxy <%s> to <%s>", t.proxyUnit.Src, t.proxyUnit.Target)
	t.server = goetty.NewServer(t.proxyUnit.Src, DECODER, ENCODER, goetty.NewInt64IdGenerator())
	t.server.Serve(t.doServe)
}

func (t *TCPServer) stop() {
	t.server.Stop()
}

func (t *TCPServer) pause() {
	t.Lock()
	if t.paused {
		t.Unlock()
		return
	}
	t.paused = true
	t.stop()
	t.Unlock()
}

func (t *TCPServer) resume() {
	t.Lock()
	if !t.paused {
		t.Unlock()
		return
	}
	t.paused = false
	go t.start()
	t.Unlock()
}

func (t *TCPServer) doServe(session goetty.IOSession) error {
	var data interface{}
	var err error

	// client connected, make a connection to target
	conn := goetty.NewConnector(t.createGoettyConf(), DECODER, ENCODER)
	_, err = conn.Connect()

	if err != nil {
		log.InfoErrorf(err, "Connect to <%s> failure.", t.proxyUnit.Target)
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
			ctl := t.proxyUnit.Ctl

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
			ctl := t.proxyUnit.Ctl

			if 0 == ctl.Out.LossRate {
				t.doWrite(bytes, conn, ctl.Out)
			} else {
				if rand.Intn(100) > ctl.Out.LossRate {
					t.doWrite(bytes, conn, ctl.Out)
				} else {
					log.Infof("Loss write <%+v> to <%s>", bytes, t.proxyUnit.Target)
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
		Addr:                   t.proxyUnit.Target,
		TimeWheel:              util.GetTimeWheel(),
		TimeoutWrite:           time.Second * time.Duration(t.proxyUnit.TimeoutWrite),
		TimeoutConnectToServer: time.Second * time.Duration(t.proxyUnit.TimeoutConnect),
		WriteTimeoutFn:         t.writeTimeout,
	}
}

func (t *TCPServer) doWrite(bytes []byte, conn *goetty.Connector, ctl *conf.CtlUnit) {
	if ctl.DelayMs > 0 {
		log.Infof("Delay <%d>ms write to <%s>", ctl.DelayMs, t.proxyUnit.Target)
		time.Sleep(time.Millisecond * time.Duration(ctl.DelayMs))
	}

	conn.Write(bytes)
	log.Infof("Write <%+v> to <%s>", bytes, t.proxyUnit.Target)
}

func (t *TCPServer) doWriteToClient(bytes []byte, session goetty.IOSession, ctl *conf.CtlUnit) {
	if ctl.DelayMs > 0 {
		log.Infof("Delay <%d>ms write to client<%s>", ctl.DelayMs, session.RemoteAddr())
		time.Sleep(time.Millisecond * time.Duration(ctl.DelayMs))
	}

	session.Write(bytes)
	log.Infof("Write <%+v> to client<%s>", bytes, session.RemoteAddr())
}
