package auth

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		obj, err := c.Cookie("user")

		if err != nil || len(obj.Value) == 0 {
			return c.JSON(http.StatusBadRequest, nil)
		}

		return next(c)
	}
}

func CheckNotLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		obj, err := c.Cookie("user")

		if err != nil || len(obj.Value) == 0 {
			return next(c)
		}

		return c.JSON(http.StatusBadRequest, nil)
	}
}

func CheckAuthAPI(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// c.Logger().Print(c.Request().Header)

		token, ok := c.Request().Header["Token"]

		// c.Logger().Print(token, ok)

		if !ok || len(token) == 0 {
			return c.JSON(http.StatusBadRequest, nil)
		}

		correct := os.Getenv("_osch_api_token")

		// c.Logger().Print(correct)

		if correct != token[0] {
			return c.JSON(http.StatusForbidden, nil)
		}

		return next(c)
	}
}
