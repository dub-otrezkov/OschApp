package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

type MySQLdatabase struct {
	db *sql.DB
}

type BadConnection struct {
	i string
}

func (b BadConnection) Error() string {
	return fmt.Sprintf("failed to connect %v", b.i)
}

func New() (*MySQLdatabase, error) {
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
		return nil, BadConnection{i: "_osch_user"}
	}
	dt.passwd, exist = os.LookupEnv("_osch_passwd")
	if !exist {
		return nil, BadConnection{i: "_osch_passwd"}
	}
	dt.addr, exist = os.LookupEnv("_osch_addr")
	if !exist {
		return nil, BadConnection{i: "_osch_addr"}
	}
	dt.bdname, exist = os.LookupEnv("_osch_bdname")
	if !exist {
		return nil, BadConnection{i: "_osch_bdname"}
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
		return nil, err
	}

	return res, nil
}

func (d *MySQLdatabase) Query(query string, args ...any) (*sql.Rows, error) {
	cn, err := d.db.Query(query, args...)
	return cn, err
}

func (d *MySQLdatabase) Exec(query string, args ...any) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

func (d *MySQLdatabase) Get_columns(dbname string) ([]string, error) {
	var rows *sql.Rows
	rows, err := (*d.db).Query(fmt.Sprintf("select * from %v", dbname))
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
