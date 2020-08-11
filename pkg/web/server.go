package web

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yfedoruck/book/pkg/env"
	"github.com/yfedoruck/book/pkg/pg"
	"html/template"
)

type Server struct {
	Port   string
	Echo   *echo.Echo
	Router *Router
	db     *pg.Postgres
}

func (s Server) Start() {
	// Echo instance
	e := s.Echo

	// Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseGlob(env.BasePath() + "/public/views/*.html")),
	}
	e.Renderer = t
	e.Static("/static", "static")
	e.Static("/book", "data")

	e.GET("/file/:id", func(c echo.Context) error {
		return c.File("data/:id")
	})

	s.Router.Routes()

	e.Logger.Fatal(e.Start(":" + env.Port()))
}

type Book struct {
	Author string
	Title  string
	Link   string
}

func NewServer(echo *echo.Echo, db *pg.Postgres) *Server {
	var s = &Server{}
	s.Port = env.Port()
	s.Echo = echo
	s.db = db
	s.Router = NewRouter(echo, db)
	return s
}
