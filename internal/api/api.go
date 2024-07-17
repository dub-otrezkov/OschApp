package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type database interface {
	Exec(query string, args ...any) (sql.Result, error)
	Get_Table(dbname string, params string) ([]map[string]interface{}, error)
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

	api.e.GET("/api/get/:dbname", api.Get_all_by_dbname)
	api.e.POST("/api/submit", api.AddSubmission)
	api.e.Static("/api/files", "files")
}

func (api *API) Get_all_by_dbname(c echo.Context) error {
	dbname := c.Param("dbname")

	api.e.Logger.Printf("%v", dbname)

	params := ""

	for k, v := range c.QueryParams() {
		if len(params) > 0 {
			params += " and "
		}
		params += fmt.Sprintf(`%v=%v`, k, v[0])
	}

	c.Logger().Print(params)

	mp, err := api.db.Get_Table(dbname, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mp)
}

func (api *API) AddSubmission(c echo.Context) error {

	type Submission struct {
		UserId string `json:"UserId"`
		TaskId string `json:"TaskId"`
		Status int    `json:"Status"`
	}

	s := Submission{}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	c.Logger().Print(body)

	for i := range body {
		c.Logger().Print(string(body[i]))
	}

	err = json.Unmarshal(body, &s)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	c.Logger().Print(s)

	_, err = api.db.Exec(fmt.Sprintf("insert into Submissions (task_id, user, status) values (%v, %v, %v)", s.TaskId, s.UserId, s.Status))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "ok", "verdict": s.Status})
}
