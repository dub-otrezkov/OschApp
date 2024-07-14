package tasks

import (
	"database/sql"
	"net/http"

	"github.com/dub-otrezkov/OschApp/pkg/auth"
	"github.com/labstack/echo"
)

type database interface {
	Get_Table(dbname string, params string) ([]map[string]interface{}, error)
	Exec(query string, args ...any) (sql.Result, error)
}

type TasksApp struct {
	e  *echo.Echo
	db database
}

func New(db database) *TasksApp {
	return &TasksApp{e: nil, db: db}
}

func (t *TasksApp) Init(e *echo.Echo) {
	t.e = e

	t.e.GET("/tasks", t.tasksList, auth.CheckLogin)
	t.e.GET("/tasks/:id", t.task, auth.CheckLogin)
	t.e.POST("/tasks/:id", t.makeSubmission, auth.CheckLogin)
}

func (*TasksApp) tasksList(c echo.Context) error {
	return c.Render(http.StatusOK, "taskslist.html", nil)
}

func (*TasksApp) task(c echo.Context) error {
	return c.Render(http.StatusOK, "task.html", struct {
		Taskid string
	}{
		Taskid: c.Param("id"),
	})
}
