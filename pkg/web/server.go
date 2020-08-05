package web

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yfedoruck/book/pkg/env"
	"html/template"
)

type Server struct {
	Port   string
	Echo   *echo.Echo
	Router *Router
}

func (s Server) Start() {
	// Echo instance
	e := s.Echo

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
	e.Static("/static", "static")
	e.Static("/book", "data")

	e.GET("/file/:id", func(c echo.Context) error {
		return c.File("data/:id")
	})

	s.Router.Routes()

	e.Logger.Fatal(e.Start(":1323"))
}

type Book struct {
	Author string
	Title  string
	Link   string
}

func NewServer(echo *echo.Echo) *Server {
	var s = &Server{}
	s.Port = env.Port()
	s.Echo = echo
	s.Router = NewRouter(echo)
	return s
}
