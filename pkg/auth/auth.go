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
	e.GET("/login", a.LoginPage, CheckNotLogin)
	e.GET("/register", a.RegisterPage, CheckNotLogin)

	e.POST("/login", a.ProcessLogin, CheckAuthAPI)
	e.POST("/register", a.ProcessRegister, CheckAuthAPI)
	e.POST("/exit", a.ProcessExit, CheckAuthAPI)
}

func (a *Auth) LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func (a *Auth) RegisterPage(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

type AuthQuery struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

func (a *Auth) ProcessLogin(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	qr := AuthQuery{}
	err = json.Unmarshal(body, &qr)

	c.Logger().Print(body)

	for _, el := range body {
		c.Logger().Print(string(el))
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	dt, err := a.db.GetTable(a.users_db_name, fmt.Sprintf(`login='%v'`, qr.Username))
	if err != nil {
		c.Logger().Print(err)
		return err
	}

	if len(dt) == 0 {
		return c.JSON(http.StatusBadRequest, "no such user")
	}

	cor := dt[0]

	if cor["password"] != qr.Password {
		return c.JSON(http.StatusBadRequest, "wrong password")
	}

	c.SetCookie(&http.Cookie{
		Name:  "user",
		Value: qr.Username,
	})

	c.SetCookie(&http.Cookie{
		Name:  "userId",
		Value: fmt.Sprint(cor["id"]),
	})

	return c.JSON(http.StatusOK, nil)
}

func (a *Auth) ProcessRegister(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "err")
	}
	qr := AuthQuery{}
	err = json.Unmarshal(body, &qr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "err")
	}

	a.db.Exec(fmt.Sprintf("insert into %v (login, password) values (%v, %v)", a.users_db_name, qr.Username, qr.Password))

	return a.ProcessLogin(c)
}

func (a *Auth) ProcessExit(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:  "user",
		Value: "",
	})

	return c.JSON(http.StatusOK, nil)
}
