package route

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/so-chiru/llct-server/dashboard"
)

func updatesHandler(w http.ResponseWriter, r *http.Request) {
	var boards = dashboard.GetDashboards()
	var notices = dashboard.GetNotices()

	var result = dashboard.LLCTUpdate{
		Updates:    time.Now().Unix(),
		Notices:    *notices,
		Dashboards: *boards,
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	CreateJsonResponse(&w, true, &bytes)
}
