package db

import (
	"fmt"
)

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
