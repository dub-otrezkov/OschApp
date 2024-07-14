package auth

import (
	"net/http"

	"github.com/labstack/echo"
)

func SetLoginHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := c.Cookie("username")
		if err == nil && len(user.Value) > 0 {
			c.Request().Header.Set("username", user.Value)
		}
		return next(c)
	}
}

func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := c.Cookie("username")

		if err != nil || len(user.Value) == 0 {
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
		user, err := c.Cookie("username")

		// c.Logger().Print(user.Name)

		if err == nil && len(user.Value) != 0 {
			t := func(c echo.Context) error {
				return c.String(http.StatusForbidden, "already authorized")
			}

			return t(c)
		}
		// c.Request().Header.Add("username", user.Value)

		return next(c)
	}
}
