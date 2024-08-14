package db

import (
	"fmt"
)

type NoUserErr struct{}

func (NoUserErr) Error() string {
	return "requested user not found"
}

func (db *MySQLdatabase) GetUser(login string) (map[string]interface{}, error) {
	t, err := db.GetTable("Users", fmt.Sprintf("login='%v'", login))
	if err != nil {
		return nil, err
	}
	if len(t) == 0 {
		return nil, NoUserErr{}
	}
	return t[0], nil
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
