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
	api.e.POST("/api/submit", api.AddSubmission, auth.CheckAuthAPI)
	api.e.POST("/api/add/:dbname", api.addObject, auth.CheckAuthAPI)

	api.e.Static("/static", "client")
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

	c.Logger().Print(params)

	mp, err := api.db.GetTable(dbname, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": err.Error()})
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

	err = json.Unmarshal(body, &s)
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

func (api *API) addObject(c echo.Context) error {
	dbname := c.Param("dbname")

	mp := make(map[string]interface{})

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": err.Error()})
	}

	err = json.Unmarshal(body, &mp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": err.Error()})
	}

	nms := ""
	vls := ""
	for key, vl := range mp {
		if len(nms) > 0 {
			nms += ", "
		}
		nms += key

		if len(vls) > 0 {
			vls += ", "
		}
		vls += fmt.Sprint(vl)
	}

	_, err = api.db.Exec("insert into %v (%v) values (%v)", dbname, nms, vls)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"status": "ok", "object": mp})
}
