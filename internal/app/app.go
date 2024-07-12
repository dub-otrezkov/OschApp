package app

import (
	"html/template"
	"log"
	"net/http"
)

type module interface {
	Init()
}

type App struct {
	port string
	md   []module
}

func New(_port string, _md ...module) *App {
	res := &App{port: _port}
	res.md = _md
	return res
}

func (a *App) Run() {
	for _, el := range a.md {
		el.Init()
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./files/general/index.html")
		if err != nil {
			panic(err)
		}

		err = t.Execute(w, nil)

		if err != nil {
			panic(err)
		}
	})

	log.Fatal(http.ListenAndServe(a.port, nil))
}
