package tasks

import (
	"net/http"

	"github.com/dub-otrezkov/OschApp/pkg/auth"
	"github.com/labstack/echo"
)

type TaskApp struct {
}

func New() *TaskApp {
	return &TaskApp{}
}

func (t *TaskApp) Init(e *echo.Echo) {
	e.GET("/tasks", t.taskslistPage, auth.CheckLogin)
}

func (*TaskApp) taskslistPage(c echo.Context) error {
	return c.Render(http.StatusOK, "taskslist.html", nil)
}
