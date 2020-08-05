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

	db := &pg.Postgres{}
	db.Connect()
	a.db = db

	a.server = web.NewServer(e)
}

func (a *App) Run() {
	defer a.db.Close()
	a.db.Tables()

	a.server.Start()
}
