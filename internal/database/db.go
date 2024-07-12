package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type MySQLdatabase struct {
	db *sql.DB
}

func New() *MySQLdatabase {
	res := &MySQLdatabase{}

	var dt struct {
		user   string
		passwd string
		addr   string
		bdname string
	}

	var exist bool
	dt.user, exist = os.LookupEnv("_osch_user")
	if !exist {
		panic("no connection 1")
	}
	dt.passwd, exist = os.LookupEnv("_osch_passwd")
	if !exist {
		panic("no connection 2")
	}
	dt.addr, exist = os.LookupEnv("_osch_addr")
	if !exist {
		panic("no connection 3")
	}
	dt.bdname, exist = os.LookupEnv("_osch_bdname")
	if !exist {
		panic("no connection 4")
	}

	cfg := mysql.Config{
		User:   dt.user,
		Passwd: dt.passwd,
		Net:    "tcp",
		Addr:   dt.addr,
		DBName: dt.bdname,
	}

	var err error
	res.db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func (d *MySQLdatabase) Query(query string, args ...any) (*sql.Rows, error) {
	cn, err := d.db.Query(query, args...)
	return cn, err
}

func (d *MySQLdatabase) Get_columns(bdname string) ([]string, error) {
	var rows *sql.Rows
	rows, err := (*d.db).Query(fmt.Sprintf("select * from %v", bdname))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []string
	cols, err = rows.Columns()
	if err != nil {
		return nil, err
	}
	return cols, nil
}
