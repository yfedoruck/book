package web

import (
	"github.com/yfedoruck/book/pkg/env"
	"github.com/yfedoruck/book/pkg/fail"
	"log"
	"net/http"
)

type Server2 struct {
	Port string
}

func (s Server2) Start() {
	log.Println("ListenAndServe port: ", s.Port)
	err := http.ListenAndServe(":" + s.Port, nil)
	fail.Check(err)
}

func NewServer2() *Server2 {
	var s = &Server2{}
	s.Port = env.Port()
	return s
}