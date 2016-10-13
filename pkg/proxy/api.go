package proxy

import (
	"net/http"

	"io/ioutil"

	"github.com/fagongzi/netproxy/pkg/conf"
	"github.com/labstack/echo"
	sd "github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	"github.com/siddontang/go/log"
)

var (
	// APIClients clients API
	APIClients = "/api/clients"
	// APIProxies proxies API
	APIProxies = "/api/proxies"
)

func (p *Proxy) startAPIServer() {
	p.apiServer.Use(mw.Logger())
	p.apiServer.Use(mw.Recover())

	p.apiServer.GET(APIProxies, p.proxies())
	p.apiServer.PUT(APIProxies, p.updateProxy())
	p.apiServer.DELETE(APIProxies, p.pause())
	p.apiServer.POST(APIProxies, p.resume())

	log.Infof("api server start at <%s>", p.cnf.APIAddr)
	p.apiServer.Run(sd.New(p.cnf.APIAddr))
}

func (p *Proxy) proxies() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, p.cnf.Units)
	}
}

func (p *Proxy) updateProxy() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctl, err := conf.UnMarshalCtlFromReader(c.Request().Body())

		if nil != err {
			return c.NoContent(http.StatusBadRequest)
		}

		p.UpdateCtl(ctl)
		return c.JSON(http.StatusOK, ctl)
	}
}

func (p *Proxy) pause() echo.HandlerFunc {
	return func(c echo.Context) error {
		addr, err := ioutil.ReadAll(c.Request().Body())
		if nil != err {
			return c.NoContent(http.StatusBadRequest)
		}
		p.Pause(string(addr))
		return c.JSON(http.StatusOK, "OK")
	}
}

func (p *Proxy) resume() echo.HandlerFunc {
	return func(c echo.Context) error {
		addr, err := ioutil.ReadAll(c.Request().Body())
		if nil != err {
			return c.NoContent(http.StatusBadRequest)
		}
		p.Resume(string(addr))
		return c.JSON(http.StatusOK, "OK")
	}
}
