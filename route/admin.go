package route

import (
	"net/http"
	"os"

	"github.com/so-chiru/llct-server/metaify/metautils"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	var pass = r.URL.Query().Get("pass")
	var PASS = os.Getenv("ADMIN_PASSPHRASE")

	if PASS == "" || pass != PASS {
		var error_string = []byte("올바르지 않은 요청입니다.")
		CreateJsonResponse(&w, false, &error_string)

		return
	}

	result, err := metautils.PullUpdate()
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)

		return
	}

	var data = []byte(result)
	CreateJsonResponse(&w, true, &data)
}
