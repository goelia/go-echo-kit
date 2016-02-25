package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/cors"
	"github.com/thoas/stats"
	"net"
	"github.com/nu7hatch/gouuid"
	"github.com/goelia/go-echo-kit/errs"
	"net/http"
)

func main() {
	e := echo.New()
	e.SetDebug(true)

	e.SetHTTPErrorHandler(func(err error, c *echo.Context) {
		code := http.StatusInternalServerError
		msg := http.StatusText(code)
		if er, ok := err.(*errs.Err); ok {
			code = 400
			msg = er.Error()
		}
		if !c.Response().Committed() {
			http.Error(c.Response(), msg, code)
		}
		e.Logger().Error(err)
	})

	// -----------
	// Middleware
	// -----------
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//------------------------
	// Third-party middleware
	//------------------------
	e.Use(cors.Default().Handler) // https://github.com/rs/cors

	s := stats.New() // https://github.com/thoas/stats
	e.Use(s.Handler)

	e.Use(func(c *echo.Context) error {
		req := c.Request()
		remoteAddr := req.RemoteAddr
		if ip := req.Header.Get(echo.XRealIP); ip != "" {
			remoteAddr = ip
		} else if ip = req.Header.Get(echo.XForwardedFor); ip != "" {
			remoteAddr = ip
		} else {
			remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
		}
		c.Set("ip", remoteAddr)
		return nil
	})

	// Map uuid for every requests
	e.Use(func(c *echo.Context) error {
		id, _ := uuid.NewV4()
		c.Set("uuid", id)
		return nil
	})

	e.Get("/", func(c *echo.Context) error {
		return c.String(200, "hello world!")
	})

	v := e.Group("/v1")
	v.Get("/routes", func(c *echo.Context) error {
		return c.JSON(200, e.Routes())
	})
	v.Get("/stats", func(c *echo.Context) error {
		return c.JSON(200, s.Data())
	})

	t := e.Group("/test")
	t.Get("/error", func(c *echo.Context) error {
		return &errs.Err{Code:errs.BadRequest}
	})
	e.Run(":3000")
}
