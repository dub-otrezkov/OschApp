package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	mdl "github.com/dub-otrezkov/OschApp/internal/database"
	"github.com/dub-otrezkov/OschApp/pkg/auth"
	"github.com/labstack/echo"
)

type database interface {
	AddSubmision(mdl.Submission) error
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
		if v[0] != "null" {
			params += fmt.Sprintf(`%v=%v`, k, v[0])
		} else {
			params += fmt.Sprintf(`%v is %v`, k, v[0])
		}
	}

	// c.Logger().Print(params)

	mp, err := api.db.GetTable(dbname, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": err.Error()})
	}

	return c.JSON(http.StatusOK, mp)
}

func (api *API) Submit(c echo.Context) error {

	type jsub struct {
		TaskId    int    `json:"TaskId"`
		SessionId int    `json:"SessionId"`
		Answer    string `json:"Answer"`
	}

	s := jsub{}

	err := getReqBody(&c, &s)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"verdict": 3})
	}

	c.Logger().Print(s)

	res := mdl.Submission{TaskId: s.TaskId, SessionId: s.SessionId, Verdict: 0}
	if cor, err := api.db.GetTable("Tasks", fmt.Sprintf("id=%v", s.TaskId)); err == nil && cor[0]["ans"] == s.Answer {
		res.Verdict = 1
	}

	err = api.db.AddSubmision(res)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]any{"verdict": 3})
	}

	return c.JSON(http.StatusOK, map[string]any{"verdict": res.Verdict})
}
