package route

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/so-chiru/llct-server/dashboard"
	"github.com/so-chiru/llct-server/structs"
)

func updatesHandler(w http.ResponseWriter, r *http.Request) {
	var boards = dashboard.GetDashboards()
	var notices = dashboard.GetNotices()

	var result = structs.LLCTUpdate{
		UpdateAt:  time.Now().Unix(),
		Notices:   *notices,
		Dashboard: *boards,
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}
