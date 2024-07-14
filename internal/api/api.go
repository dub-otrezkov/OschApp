package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type database interface {
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
