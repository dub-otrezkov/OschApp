package userstats

import (
	"net/http"

	"github.com/dub-otrezkov/OschApp/pkg/auth"
	"github.com/labstack/echo"
)

type Userstats struct {
}

func New() *Userstats {
	return &Userstats{}
}

func (u *Userstats) Init(e *echo.Echo) {
	e.GET("/stats", u.ExamResults, auth.CheckLogin)
}

func (u *Userstats) ExamResults(c echo.Context) error {
	return c.Render(http.StatusOK, "userstats.html", nil)
}
