package route

import "net/http"

func updatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update"))
}
