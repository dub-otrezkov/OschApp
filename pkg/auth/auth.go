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
	GetUser(login string) ([]map[string]interface{}, error)
	RegisterUser(login string, password string) error
}

type Auth struct {
	db database
}

func New(db database) *Auth {
	return &Auth{db: db}
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
		c.Logger().Print(err.Error())

		return c.JSON(http.StatusBadRequest, nil)
	}
	qr := AuthQuery{}
	err = json.Unmarshal(body, &qr)

	c.Logger().Print(body)

	for _, el := range body {
		c.Logger().Print(string(el))
	}

	if err != nil {
		c.Logger().Print(err.Error())

		return c.JSON(http.StatusBadRequest, nil)
	}

	// dt, err := a.db.GetTable(a.users_db_name, fmt.Sprintf(`login='%v'`, qr.Username))
	dt, err := a.db.GetUser(qr.Username)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
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

	// a.db.Exec(fmt.Sprintf("insert into %v (login, password) values (%v, %v)", a.users_db_name, qr.Username, qr.Password))
	err = a.db.RegisterUser(qr.Username, qr.Password)

	if err == nil {
		return a.ProcessLogin(c)
	}

	return c.JSON(http.StatusBadRequest, err.Error())
}

func (a *Auth) ProcessExit(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:  "user",
		Value: "",
	})

	return c.JSON(http.StatusOK, nil)
}
