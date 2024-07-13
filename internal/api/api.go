package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type database interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Get_columns(bdname string) ([]string, error)
}

type API struct {
	db *database
	e  *echo.Echo
}

func New(db database) *API {
	return &API{db: &db, e: nil}
}

func (api *API) Init(e *echo.Echo) {
	api.e = e

	api.e.GET("/api/get/:dbname", api.Get_all_by_dbname)
	api.e.Static("/api/files", "files")
}

func (api *API) Get_all_by_dbname(c echo.Context) error {
	dbname := c.Param("dbname")

	api.e.Logger.Printf("%v", dbname)

	cols, err := (*api.db).Get_columns(dbname)
	if err != nil {
		return err
	}

	prm := ""

	for i := range cols {
		vl := c.QueryParam(cols[i])
		if len(vl) == 0 {
			continue
		}
		if prm != "" {
			prm += " and "
		}
		prm += fmt.Sprintf("%v='%v'", cols[i], vl)
	}

	qry := ""
	if len(prm) == 0 {
		qry = fmt.Sprintf("select * from %v", dbname)
	} else {
		qry = fmt.Sprintf("select * from %v where %v", dbname, prm)
	}

	log.Println(qry)

	cn, err := (*api.db).Query(qry)
	if err != nil {
		return err
	}

	defer cn.Close()

	var mp []map[string]interface{}

	cont := make([]interface{}, len(cols))
	var lks = make([]interface{}, len(cols))
	for i := range cont {
		lks[i] = &cont[i]
	}

	for cn.Next() {
		err = cn.Scan(lks...)

		if err != nil {
			return err
		}

		var tmp = make(map[string]interface{})

		for j := 0; j < len(cols); j++ {
			switch cont[j].(type) {
			case []byte:
				tmp[cols[j]] = string(cont[j].([]byte))
			default:
				tmp[cols[j]] = cont[j]
			}
		}

		mp = append(mp, tmp)
	}

	api.e.Logger.Print(mp)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mp)
}
