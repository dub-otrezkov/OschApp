package auth

import (
	"net/http"

	"github.com/labstack/echo"
)

type Auth struct {
}

func New() *Auth {
	return &Auth{}
}

func (a *Auth) Init(e *echo.Echo) {
	e.GET("/login", a.LoginPage, CheckNotLogin)
	e.POST("/login", a.ProcessLogin, CheckNotLogin)
	e.POST("/exit", a.Exit, CheckLogin)
}

func (Auth) LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func (Auth) ProcessLogin(c echo.Context) error {
	c.Logger().Print(c.Request().FormValue("username"))

	c.SetCookie(&http.Cookie{Name: "username", Value: c.Request().FormValue("username")})
	http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusOK)
	return nil
}

func (Auth) Exit(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: "username", Value: ""})
	http.Redirect(c.Response().Writer, c.Request(), "/", http.StatusOK)
	return nil
}
