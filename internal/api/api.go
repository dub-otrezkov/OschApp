package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	mdl "github.com/dub-otrezkov/OschApp/internal/database"
	"github.com/labstack/echo"
)

type database interface {
	AddSubmision(mdl.Submission) error
	CloseSession(session_id int) error
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

	// endpoints
	api.e.GET("/api/get/:dbname", api.getTable)
	api.e.GET("/api/stats/:user_id", api.getUserStats)

	api.e.POST("/api/submit", api.submit)
	api.e.POST("/api/finish", api.finishExam)

	// files
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

func (api *API) getTable(c echo.Context) error {
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

	mp, err := api.db.GetTable(dbname, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": err.Error()})
	}

	return c.JSON(http.StatusOK, mp)
}
