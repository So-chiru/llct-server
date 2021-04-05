package route

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/so-chiru/llct-server/utils"
)

func openCoverFile(group_name string, id string) (bool, []byte) {
	var base_folder = "./datas/" + group_name + "/" + id + "/"

	if _, err := os.Stat(base_folder); os.IsNotExist(err) {
		return false, nil
	}

	reader, err := os.Open(base_folder + "cover.png")
	if err != nil {
		return false, nil
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return false, nil
	}

	return true, bytes
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

	song_number, err := strconv.Atoi(id[1:])
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	exists, bytes := openCoverFile(group_name, fmt.Sprint(song_number))

	if !exists {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(bytes)
}
