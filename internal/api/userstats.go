package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func (api *API) getUserStats(c echo.Context) error {
	u := c.Param("user_id")

	sessions, err := api.db.GetTable("Sessions", fmt.Sprintf("id>0 and user_id=%v and active=0", u))

	if err != nil {
		c.Logger().Print(err.Error())
		return c.JSON(http.StatusBadRequest, make([]int, 0))
	}

	res := make([]map[string]interface{}, len(sessions))

	for i, el := range sessions {
		c.Logger().Print(el)

		res[i] = make(map[string]interface{})

		res[i]["exam_id"] = int(el["exam_id"].(int64))

		tsks, err := api.db.GetTable("Tasklist", fmt.Sprintf("exam_id=%v", res[i]["exam_id"]))
		if err != nil {
			c.Logger().Print(err.Error())
			return c.JSON(http.StatusBadRequest, make([]int, 0))
		}

		subs, err := api.db.GetTable("Submissions", fmt.Sprintf("session_id=%v", el["id"]))

		if err != nil {
			c.Logger().Print(err.Error())
			return c.JSON(http.StatusBadRequest, make([]int, 0))
		}

		st := make(map[int64]int64)

		for _, j := range subs {
			st[j["task_id"].(int64)] = max(st[j["task_id"].(int64)], j["status"].(int64))

		}

		c.Logger().Print(tsks)

		ans := make([]map[string]any, 0, len(tsks))

		for _, j := range tsks {
			vl, ok := st[j["task_id"].(int64)]
			t1 := -1
			if ok {
				t1 = int(vl)
			}

			ans = append(ans, map[string]any{"id": j["task_id"].(int64), "status": t1})

		}

		res[i]["ans"] = ans
	}

	return c.JSON(http.StatusOK, res)

}
