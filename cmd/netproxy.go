package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"encoding/json"
	"io/ioutil"

	"github.com/CodisLabs/codis/pkg/utils/log"
	"github.com/fagongzi/netproxy/pkg/conf"
	l "github.com/fagongzi/netproxy/pkg/log"
	"github.com/fagongzi/netproxy/pkg/proxy"
	"github.com/fagongzi/netproxy/pkg/util"
)

var (
	cpus     = flag.Int("cpus", 1, "use cpu nums")
	file     = flag.String("config", "", "config file")
	logFile  = flag.String("log-file", "", "which file to record log, if not set stdout to use.")
	logLevel = flag.String("log-level", "info", "log level.")
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*cpus)

	l.InitLog(*logFile)
	l.SetLogLevel(*logLevel)

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		log.PanicErrorf(err, "read config file <%s> failure.", *file)
	}

	cnf := &conf.Conf{}
	err = json.Unmarshal(data, cnf)
	if err != nil {
		log.PanicErrorf(err, "parse config file <%s> failure.", *file)
	}

	util.Init()

	p := proxy.NewProxy(cnf)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, os.Kill)

	go func() {
		<-c
		log.Infof("ctrl-c or SIGTERM found, netproxy will exit")
		p.Stop()
	}()

	p.Start()
	log.Infof("netproxy is Exit.")
}
