package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type database interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Get_columns(bdname string) ([]string, error)
}

type API struct {
	db *database
}

func New() *API {
	return &API{}
}

func (api *API) Get_all_by_dbname(w http.ResponseWriter, r *http.Request) {
	bdname := r.PathValue("bdname")

	cols, err := (*api.db).Get_columns(bdname)
	if err != nil {
		panic(err)
	}

	queryParams := r.URL.Query()

	prm := ""

	for i := range cols {
		vl := queryParams.Get(cols[i])
		if len(vl) == 0 {
			continue
		}
		if prm != "" {
			prm += " and "
		}
		prm += fmt.Sprintf("%v='%v'", cols[i], vl)
	}

	// fmt.Printf("cols: %v %v\n", cols, len(cols))

	qry := ""
	if len(prm) == 0 {
		qry = fmt.Sprintf("select * from %v", bdname)
	} else {
		qry = fmt.Sprintf("select * from %v where %v", bdname, prm)
	}

	log.Println(qry)

	cn, err := (*api.db).Query(qry)
	if err != nil {
		panic(err)
	}

	defer cn.Close()

	var mp []map[string]interface{}

	// return
	cont := make([]interface{}, len(cols))
	var lks = make([]interface{}, len(cols))
	for i := range cont {
		lks[i] = &cont[i]
	}

	for cn.Next() {
		err = cn.Scan(lks...)

		if err != nil {
			panic(err)
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

	log.Printf("mp: %v\n", mp)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Charset", "utf-8")
	json.NewEncoder(w).Encode(mp)
}
