package tasks

import (
	"net/http"

	"github.com/dub-otrezkov/OschApp/pkg/auth"
	"github.com/labstack/echo"
)

type TasksApp struct {
	e *echo.Echo
}

func New() *TasksApp {
	return &TasksApp{e: nil}
}

func (t *TasksApp) Init(e *echo.Echo) {
	t.e = e

	t.e.GET("/tasks", t.tasks_list_page, auth.CheckLogin)
	t.e.GET("/tasks/:id", t.task_page, auth.CheckLogin)
}

func (TasksApp) tasks_list_page(c echo.Context) error {
	return c.Render(http.StatusOK, "taskslist.html", nil)
}

func (TasksApp) task_page(c echo.Context) error {
	return c.Render(http.StatusOK, "task.html", struct {
		Taskid string
	}{
		Taskid: c.Param("id"),
	})
}
