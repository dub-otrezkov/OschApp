package app

import (
	"log"
	"net/http"
)

type handler interface {
	init()
}

type App struct {
	port string
	md   []handler
}

func New(_port string, _md ...handler) *App {
	res := &App{port: _port}
	res.md = _md
	return res
}

func (a *App) Run() {
	for _, el := range a.md {
		el.init()
	}
	log.Fatal(http.ListenAndServe(a.port, nil))
}
