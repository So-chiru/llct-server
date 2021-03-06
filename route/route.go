package route

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func emptyResponse(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Couldn't create a hijacker.", http.StatusInternalServerError)
		return
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn.Close()
}

func CreateJsonResponse(w *http.ResponseWriter, s bool, d *[]byte) {
	var result_string string

	if s {
		result_string = "success"
	} else {
		result_string = "error"
	}

	var data interface{}
	err := json.Unmarshal(*d, &data)
	if err != nil {
		data = string(*d)
	}

	res := JsonResponse{
		Result: result_string,
		Data:   data,
	}

	encoder := json.NewEncoder(*w)

	encoder.Encode(res)
}

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", emptyResponse)
	r.HandleFunc("/lists", listsHandler)
	r.HandleFunc("/updates", updatesHandler)

	r.HandleFunc("/cover/{id:[0-9]+}", coverHandler)
	r.HandleFunc("/color/{id:[0-9]+}", colorHandler)
	r.HandleFunc("/audio/{id:[0-9]+}", audioHandler)
	r.HandleFunc("/call/{id:[0-9]+}", callHandler)
	r.HandleFunc("/admin", adminHandler)

	r.HandleFunc("/spotify/oauth", spotifyLoginHandler)
	r.HandleFunc("/spotify/callback", spotifyCallbackHandler)
	r.HandleFunc("/spotify/refresh_token", spotifyRefreshTokenHandler)

	return r
}
