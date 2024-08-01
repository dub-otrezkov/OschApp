package tasks

import (
	"fmt"
	"net/http"
	"strconv"

	mdl "github.com/dub-otrezkov/OschApp/internal/database"
	"github.com/dub-otrezkov/OschApp/pkg/auth"
	"github.com/labstack/echo"
)

type database interface {
	GetTable(dbname string, params string) ([]map[string]interface{}, error)
	AddSession(mdl.Session) (int, error)
}

type TaskApp struct {
	db database
}

func New(db database) *TaskApp {
	return &TaskApp{db: db}
}

func (t *TaskApp) Init(e *echo.Echo) {
	e.GET("/tasks", t.taskslistPage, auth.CheckLogin)
	e.GET("/tasks/:id", t.taskPage, auth.CheckLogin)

	e.GET("/exams", t.examslistPage, auth.CheckLogin)
	e.GET("/exams/:id", t.examPage, auth.CheckLogin)
}

func (*TaskApp) taskslistPage(c echo.Context) error {
	return c.Render(http.StatusOK, "taskslist.html", nil)
}

func (*TaskApp) taskPage(c echo.Context) error {
	u, _ := c.Cookie("userId")
	c.SetCookie(&http.Cookie{Name: "session", Value: fmt.Sprintf("-%v", u.Value)})

	return c.Render(http.StatusOK, "task.html", struct {
		Id string
	}{
		Id: c.Param("id"),
	})
}

func (*TaskApp) examslistPage(c echo.Context) error {
	return c.Render(http.StatusOK, "examslist.html", nil)
}

func (t *TaskApp) examPage(c echo.Context) error {

	x, _ := c.Cookie("userId")
	userId, _ := strconv.Atoi(x.Value)
	examId, _ := strconv.Atoi(c.Param("id"))

	cur, err := t.db.GetTable("Sessions", fmt.Sprintf("user_id=%v and active=1 and id>0 and exam_id=%v", userId, examId))

	if err != nil {
		c.Logger().Print(err)
	}

	cur_session := 0
	if len(cur) == 0 {
		cur_session, err = t.db.AddSession(mdl.Session{Id: 0, UserId: userId, Active: true, ExamId: examId})
		if err != nil {
			c.Logger().Print(err)
		}
	} else {
		cur_session, _ = strconv.Atoi(fmt.Sprint(cur[0]["id"]))
		if err != nil {
			c.Logger().Print(err)
		}
	}

	c.SetCookie(&http.Cookie{Name: "session", Value: fmt.Sprint(cur_session)})

	return c.Render(http.StatusOK, "exam.html", struct {
		Id string
	}{
		Id: c.Param("id"),
	})
}
