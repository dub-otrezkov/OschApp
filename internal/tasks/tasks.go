package tasks

import (
	"net/http"

	"github.com/labstack/echo"
)

type TaskApp struct {
}

func (t *TaskApp) Init(e *echo.Echo) {
	e.GET("tasks/", t.taskslistPage)
}

func (*TaskApp) taskslistPage(c echo.Context) error {
	return c.Render(http.StatusOK, "taskslist.html", nil)
}
