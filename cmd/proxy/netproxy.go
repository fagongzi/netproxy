package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/fagongzi/log"
	"github.com/fagongzi/netproxy/pkg/conf"
	"github.com/fagongzi/netproxy/pkg/proxy"
)

var (
	cpus = flag.Int("cpus", 1, "use cpu nums")
	file = flag.String("cfg", "", "config file")
)

func main() {
	flag.Parse()
	log.InitLog()

	runtime.GOMAXPROCS(*cpus)

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("read config file <%s> failure. err:%+v", *file, err)
	}

	cnf := &conf.Conf{}
	err = json.Unmarshal(data, cnf)
	if err != nil {
		log.Fatalf("parse config file <%s> failure. error:%+v", *file, err)
	}

	p := proxy.NewProxy(cnf)
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-c
		log.Infof("ctrl-c or SIGTERM found, netproxy will exit")
		p.Stop()
	}()

	p.Start()
	log.Infof("netproxy is Exit.")
}
