package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type database interface {
	Exec(query string, args ...any) (sql.Result, error)
	GetTable(dbname string, qry string) ([]map[string]interface{}, error)
}

type Auth struct {
	users_db_name string
	db            database
}

func New(db database, user_db_name string) *Auth {
	return &Auth{db: db, users_db_name: user_db_name}
}

func (a *Auth) Init(e *echo.Echo) {
	e.POST("/login", a.ProcessLogin, CheckAuthAPI)
	e.POST("/register", a.ProcessRegister, CheckAuthAPI)
}

type AuthQuery struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *Auth) ProcessLogin(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	qr := AuthQuery{}
	err = json.Unmarshal(body, &qr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	dt, err := a.db.GetTable(a.users_db_name, fmt.Sprintf(`login='%v'`, qr.Username))
	if err != nil {
		c.Logger().Print(err)
		return err
	}

	if len(dt) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	cor := dt[0]

	if cor["password"] != qr.Password {
		return c.JSON(http.StatusBadRequest, nil)
	}

	c.SetCookie(&http.Cookie{
		Name:  "user",
		Value: qr.Username,
	})

	return c.JSON(http.StatusOK, nil)
}

func (a *Auth) ProcessRegister(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	qr := AuthQuery{}
	err = json.Unmarshal(body, &qr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	a.db.Exec(fmt.Sprintf("insert into %v (login, password) values (%v, %v)", a.users_db_name, qr.Username, qr.Password))

	return a.ProcessLogin(c)
}
