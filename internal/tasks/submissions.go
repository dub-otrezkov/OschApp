package tasks

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func (t *TasksApp) makeSubmission(c echo.Context) error {
	user_id, err := c.Cookie("user_id")
	if err != nil {
		return err
	}

	task_id := c.Param("id")

	dt, err := t.db.Get_Table("Tasks", fmt.Sprintf("id=%v", task_id))
	if err != nil {
		c.Logger().Print(err)
		return err
	}
	if len(dt) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	task := dt[0]

	uans := c.FormValue("ans")

	c.Logger().Print(task["ans"], uans)

	var status int
	if task["ans"] == uans {
		status = 1
	} else {
		status = 0
	}

	_, err = t.db.Exec(fmt.Sprintf("insert into Submissions (user, task_id, status) values (%v, %v, %v)", user_id.Value, task_id, status))

	if err != nil {
		c.Logger().Print(err)
		return err
	}

	return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/tasks/%v", task_id))
}
