package api

import (
	"fmt"
	"net/http"

	mdl "github.com/dub-otrezkov/OschApp/internal/database"
	"github.com/labstack/echo"
)

func (api *API) submit(c echo.Context) error {

	type jsub struct {
		TaskId    int    `json:"TaskId"`
		SessionId int    `json:"SessionId"`
		Answer    string `json:"Answer"`
	}

	s := jsub{}

	err := getReqBody(&c, &s)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"verdict": 3})
	}

	c.Logger().Print(s)

	res := mdl.Submission{TaskId: s.TaskId, SessionId: s.SessionId, Verdict: 0}
	if cor, err := api.db.GetTable("Tasks", fmt.Sprintf("id=%v", s.TaskId)); err == nil && cor[0]["ans"] == s.Answer {
		res.Verdict = 1
	}

	err = api.db.AddSubmision(res)

	if err != nil {
		c.Logger().Print(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]any{"verdict": 3})
	}

	if s.SessionId > 0 {
		res.Verdict = 2
	}

	return c.JSON(http.StatusOK, map[string]any{"verdict": res.Verdict})
}

func (api *API) finishExam(c echo.Context) error {
	type req struct {
		SessionId int `json:"SessionId"`
	}

	r := req{}

	err := getReqBody(&c, &r)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	api.db.CloseSession(r.SessionId)

	return c.JSON(http.StatusOK, nil)
}
