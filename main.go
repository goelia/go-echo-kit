package main

import (
	"net"
	"net/http"

	"github.com/goelia/go-echo-kit/errs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nu7hatch/gouuid"
	"github.com/rs/cors"
	"github.com/thoas/stats"

	"strconv"

	"github.com/goelia/go-echo-kit/handles"
	"github.com/goelia/go-echo-kit/middlewares"
"github.com/goelia/go-echo-kit/config"
)

func main() {
	e := echo.New()
	e.SetDebug(true)

	e.SetHTTPErrorHandler(func(err error, c *echo.Context) {
		code := http.StatusInternalServerError
		msg := http.StatusText(code)
		if er, ok := err.(*errs.Err); ok {
			code = 422
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
		//panic(errors.New("test error"))
		return &errs.Err{Code: errs.BadRequest}
	})

	//用户登录
	v.Post("/auth/signin", handles.Signin)
	//发送验证码
	v.Post("/auth/code", handles.RefreshCode)

	auth := v.Group("/")
	auth.Use(middlewares.JWTAuth(config.GetConfig().SigningKey))
	auth.Get("auth", func(c *echo.Context) error {
		return c.JSON(200, "success auth.")
	})

	e.Run(":" + strconv.Itoa(config.GetConfig().Port))
}
