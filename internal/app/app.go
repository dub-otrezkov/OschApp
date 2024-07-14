package app

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type module interface {
	Init(*echo.Echo)
}

type App struct {
	port string
	md   []module
}

type TR struct {
	templates *template.Template
}

func (t *TR) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func New(_port string, _md ...module) *App {
	res := &App{port: _port}
	res.md = _md
	return res
}

func (a *App) Run() {
	e := echo.New()
	e.Renderer = &TR{templates: template.Must(template.ParseGlob("./files/html/*.html"))}

	for _, el := range a.md {
		el.Init(e)
	}

	e.GET("/", func(c echo.Context) error {
		e.Logger.Printf("main page\n")

		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.Logger.Fatal(e.Start(a.port))
}
