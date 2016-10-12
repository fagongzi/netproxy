package proxy

import (
	"net/http"

	"io/ioutil"

	"github.com/labstack/echo"
	sd "github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	"github.com/siddontang/go/log"
)

func (p *Proxy) startAPIServer() {
	p.apiServer.Use(mw.Logger())
	p.apiServer.Use(mw.Recover())

	p.apiServer.GET("/api/clients", p.clients())
	p.apiServer.PUT("/api/clients", p.updateRate())
	p.apiServer.PUT("/api/proxy/pause", p.pause())
	p.apiServer.PUT("/api/proxy/resume", p.resume())

	log.Infof("api server start at <%s>", p.cnf.APIAddr)
	p.apiServer.Run(sd.New(p.cnf.APIAddr))
}

func (p *Proxy) clients() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, p.GetAllClients())
	}
}

func (p *Proxy) updateRate() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctl, err := UnMarshalCtlFromReader(c.Request().Body())

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
