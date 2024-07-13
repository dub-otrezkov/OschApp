package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{Name: "zvuk", Value: "petr"})
		return next(c)
	}
}
