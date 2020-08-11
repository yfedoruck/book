package main

import (
	"github.com/labstack/echo/v4"
	"github.com/yfedoruck/book/pkg/pg"
	"github.com/yfedoruck/book/pkg/web"
)

type App struct {
	server    *web.Server
	framework *echo.Echo
	db        *pg.Postgres
}

func (a *App) Init() {
	e := echo.New()
	a.framework = e
	a.framework.Debug = false

	db := pg.NewPostgres()
	db.Connect()
	a.db = db

	a.server = web.NewServer(a.framework, a.db)
}

func (a *App) Run() {
	defer a.db.Close()
	a.db.CreateTables()

	a.server.Start()
}
