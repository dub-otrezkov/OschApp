package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	hsh "github.com/dub-otrezkov/OschApp/pkg/hasher"
	"github.com/labstack/echo"
)

type database interface {
	GetUser(login string) (map[string]interface{}, error)
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

	e.POST("/login", a.ProcessLogin)
	e.POST("/register", a.ProcessRegister)
	e.POST("/exit", a.ProcessExit)
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

	c.Logger().Print(body)

	qr := AuthQuery{}
	err = json.Unmarshal(body, &qr)

	if err != nil {
		c.Logger().Print(err.Error())

		return c.JSON(http.StatusBadRequest, nil)
	}

	c.Logger().Print(qr)

	qr.Password = hsh.CalcSha256(qr.Password)

	cor, err := a.db.GetUser(qr.Username)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if cor["password"] != qr.Password {
		return c.JSON(http.StatusBadRequest, "wrong password")
	}

	c.SetCookie(&http.Cookie{
		Name:  "user",
		Value: qr.Username,
	})

	c.SetCookie(&http.Cookie{
		Name:  "user_id",
		Value: fmt.Sprint(cor["id"]),
	})

	return c.JSON(http.StatusOK, nil)
}

func (a *Auth) ProcessRegister(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	c.Request().Body = io.NopCloser(bytes.NewReader(body))

	qr := AuthQuery{}
	err = json.Unmarshal(body, &qr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	qr.Password = hsh.CalcSha256(qr.Password)

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
