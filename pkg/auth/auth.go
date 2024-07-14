package auth

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type database interface {
	Exec(query string, args ...any) (sql.Result, error)
	Get_Table(dbname string, qry string) ([]map[string]interface{}, error)
}

type Auth struct {
	users_db_name string
	db            database
}

type AuthErr struct {
	err string
}

func (err AuthErr) Error() string {
	return err.err
}

func New(db database, user_db_name string) *Auth {
	return &Auth{db: db, users_db_name: user_db_name}
}

func (a *Auth) Init(e *echo.Echo) {
	e.GET("/login", a.LoginPage, CheckNotLogin)
	e.GET("/register", a.RegisterPage, CheckNotLogin)
	e.POST("/login", a.ProcessLogin, CheckNotLogin)
	e.POST("/register", a.ProcessRegister, CheckNotLogin)
	e.POST("/exit", a.Exit, CheckLogin)
}

func (*Auth) LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func (*Auth) RegisterPage(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

func (a *Auth) ProcessLogin(c echo.Context) error {
	c.Logger().Print(c.Request().FormValue("username"))

	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")

	dt, err := a.db.Get_Table(a.users_db_name, fmt.Sprintf(`login='%v'`, username))
	if err != nil {
		c.Logger().Print(err)
		return err
	}

	if len(dt) == 0 {
		c.Logger().Print(err)
		return AuthErr{err: "no such user"}
	}

	cor := dt[0]

	if cor["password"] != password {
		c.Logger().Print(password, cor["password"].(string))
		return AuthErr{err: "wrong password"}
	}

	c.SetCookie(&http.Cookie{Name: "username", Value: c.Request().FormValue("username")})
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func (a *Auth) ProcessRegister(c echo.Context) error {
	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")

	c.Logger().Print(username, password)

	a.db.Exec(fmt.Sprintf("insert into %v (login, password) values (%v, %v)", a.users_db_name, username, password))

	return a.ProcessLogin(c)
}

func (Auth) Exit(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: "username", Value: ""})
	return c.Redirect(http.StatusMovedPermanently, "/")
}
