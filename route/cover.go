package route

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/so-chiru/llct-server/utils"
)

func getCoverFilePath(group_name string, id string) string {
	var base_folder = "./datas/" + group_name + "/" + id + "/"

	if _, err := os.Stat(base_folder); os.IsNotExist(err) {
		return ""
	}

	if _, err := os.Stat(base_folder + "cover.jpg"); !os.IsNotExist(err) {
		return base_folder + "cover.jpg"
	}

	if _, err := os.Stat(base_folder + "cover.png"); !os.IsNotExist(err) {
		return base_folder + "cover.png"
	}

	return ""
}

func coverHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if len(vars["id"]) < 1 {
		var error_string = []byte("경로에 :id를 포함하세요.")
		CreateJsonResponse(&w, false, &error_string)

		return
	}

	var id = vars["id"]

	group_number, err := strconv.Atoi(string(id[0]))
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	exists, group_name := utils.GetGroupID(group_number)

	if !exists {
		var error_string = []byte("해당하는 그룹의 값이 없습니다.")
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	if len(id[1:]) > 5 {
		var error_string = []byte("올바르지 않은 ID 값입니다.")
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	song_number, err := strconv.Atoi(id[1:])
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	path := getCoverFilePath(group_name, fmt.Sprint(song_number))

	if len(path) < 1 {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	cacheFile(w, r, path)
}
