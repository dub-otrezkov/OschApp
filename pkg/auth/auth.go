package auth

import (
	"net/http"

	"github.com/labstack/echo"
)

func SetLoginHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := c.Cookie("username")
		if err == nil {
			c.Request().Header.Set("username", user.Value)
		}
		return next(c)
	}
}

func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := c.Cookie("username")

		if err != nil {
			t := func(c echo.Context) error {
				return c.String(http.StatusUnauthorized, "auth required")
			}

			return t(c)
		}
		return next(c)
	}
}

func CheckNotLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := c.Cookie("username")

		if err == nil {
			t := func(c echo.Context) error {
				return c.String(http.StatusForbidden, "already authorized")
			}

			return t(c)
		}
		// c.Request().Header.Add("username", user.Value)

		return next(c)
	}
}

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
