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
		res[i] = make(map[string]interface{})

		res[i]["exam_id"] = int(el["exam_id"].(int64))

		ans := make(map[int64]int)

		tsks, err := api.db.GetTable("Tasklist", fmt.Sprintf("examId=%v", res[i]["exam_id"]))
		if err != nil {
			c.Logger().Print(err.Error())
			return c.JSON(http.StatusBadRequest, make([]int, 0))
		}

		for _, j := range tsks {
			ans[j["taskId"].(int64)] = 2
		}

		subs, err := api.db.GetTable("Submissions", fmt.Sprintf("session_id=%v", el["id"]))
		c.Logger().Print(subs, el["id"])

		if err != nil {
			c.Logger().Print(err.Error())
			return c.JSON(http.StatusBadRequest, make([]int, 0))
		}

		for _, j := range subs {
			// c.Logger().Print(j["status"] == 0)

			if j["status"].(int64) == 0 && ans[j["task_id"].(int64)] == 2 {
				ans[j["task_id"].(int64)] = 0
			} else if j["status"].(int64) == 1 {
				ans[j["task_id"].(int64)] = 1
			}
		}

		res[i]["ans"] = ans
	}

	return c.JSON(http.StatusOK, res)

}
