package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dub-otrezkov/OschApp/pkg/auth"
	"github.com/labstack/echo"
)

type database interface {
	Exec(query string, args ...any) (sql.Result, error)
	GetTable(dbname string, params string) ([]map[string]interface{}, error)
}

type API struct {
	db database
	e  *echo.Echo
}

func New(db database) *API {
	return &API{db: db, e: nil}
}

func (api *API) Init(e *echo.Echo) {
	api.e = e

	api.e.GET("/api/get/:dbname", api.GetTable, auth.CheckAuthAPI)
	api.e.POST("/api/submit", api.Submit, auth.CheckAuthAPI)

	api.e.Static("/static", "client")
	api.e.Static("/files", "files")
}

func getReqBody(c *echo.Context, r any) error {
	body, err := io.ReadAll((*c).Request().Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, r)

	return err
}

func (api *API) GetTable(c echo.Context) error {
	dbname := c.Param("dbname")

	api.e.Logger.Printf("%v", dbname)

	params := ""

	for k, v := range c.QueryParams() {
		if len(params) > 0 {
			params += " and "
		}
		params += fmt.Sprintf(`%v=%v`, k, v[0])
	}

	// c.Logger().Print(params)

	mp, err := api.db.GetTable(dbname, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": err.Error()})
	}

	return c.JSON(http.StatusOK, mp)
}

func (api *API) Submit(c echo.Context) error {

	s := Submission{}

	err := getReqBody(&c, &s)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"verdict": 3})
	}

	verdict := 0
	if cor, err := api.db.GetTable("Tasks", fmt.Sprintf("id=%v", s.TaskId)); err == nil && cor[0]["ans"] == s.Answer {
		verdict = 1
	}

	// c.Logger().Print(fmt.Sprintf("insert into Submissions (task_id, user, status) values (%v, %v, %v)", s.TaskId, s.UserId, verdict))

	_, err = api.db.Exec(fmt.Sprintf("insert into Submissions (task_id, session_id, status) values (%v, %v, %v)", s.TaskId, s.SessionId, verdict))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"verdict": 3})
	}

	return c.JSON(http.StatusOK, map[string]any{"verdict": verdict})
}
