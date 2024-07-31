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

func (d *MySQLdatabase) GetTable(dbname string, params string) ([]map[string]interface{}, error) {
	qry := fmt.Sprintf("select * from %v where %v", dbname, params)
	if len(params) == 0 {
		qry = fmt.Sprintf("select * from %v", dbname)
	}

	fmt.Println(qry)

	cn, err := d.db.Query(qry)
	if err != nil {
		return nil, err
	}
	defer cn.Close()

	var mp []map[string]interface{}

	cols, err := cn.Columns()
	if err != nil {
		return nil, err
	}

	cont := make([]interface{}, len(cols))
	var lks = make([]interface{}, len(cols))
	for i := range cont {
		lks[i] = &cont[i]
	}

	for cn.Next() {
		err = cn.Scan(lks...)

		if err != nil {
			return nil, err
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
	return mp, nil
}

func (db *MySQLdatabase) AddSubmision(s Submission) error {
	_, err := db.db.Exec(fmt.Sprintf("insert into Submissions (task_id, session_id, status) values (%v, %v, %v)", s.TaskId, s.SessionId, s.Verdict))
	return err
}

func (db *MySQLdatabase) AddSession(s Session) (int, error) {
	var err error
	var res sql.Result
	if s.Id != 0 {
		res, err = db.db.Exec(fmt.Sprintf("insert into Sessions (id, user_id, active, exam_id) values (%v, %v, %v, null)", s.Id, s.UserId, s.Active))
	} else {
		res, err = db.db.Exec(fmt.Sprintf("insert into Sessions (user_id, active, exam_id) values (%v, %v, %v)", s.UserId, s.Active, s.ExamId))
	}

	if err != nil {
		return 0, err
	}

	r, err := res.LastInsertId()

	return int(r), err
}

func (db *MySQLdatabase) GetUser(login string) ([]map[string]interface{}, error) {
	return db.GetTable("Users", fmt.Sprintf("login='%v'", login))
}

func (db *MySQLdatabase) RegisterUser(login string, password string) error {
	res, err := db.Exec(fmt.Sprintf("insert into Users (login, password) values ('%v', '%v')", login, password))

	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(id)

	_, err = db.AddSession(Session{Id: -int(id), UserId: int(id), Active: true})

	return err
}
