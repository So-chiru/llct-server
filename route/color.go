package route

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/so-chiru/llct-server/pastelize"
	"github.com/so-chiru/llct-server/utils"
)

func checkColorFileExists(group_name string, id string) (bool, string) {
	var path = "./datas/" + group_name + "/" + id + "/" + "_cache/color.json"
	return isFileExists(path), path
}

func colorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Add("Access-Control-Allow-Headers", "llct-api-version")

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

	var path = "./datas/" + group_name + "/" + fmt.Sprint(song_number) + "/cover.jpg"

	if !isFileExists(path) {
		var error_string = []byte("커버 파일이 없습니다.")
		CreateJsonResponse(&w, false, &error_string)

		return
	}

	file, err := os.Open(path)
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	image, _, err := image.Decode(file)
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	var api_version = r.Header.Get("llct-api-version")

	if len(api_version) < 1 || api_version != "2" {
		var error_string = []byte("API 버전을 지정하지 않았거나 더 이상 지원하지 않는 API 버전입니다. 페이지를 업데이트 해주세요.")
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	cache_exists, path := checkColorFileExists(group_name, fmt.Sprint(song_number))
	if cache_exists {
		file, err := os.Open(path)
		if err != nil {
			var error_string = []byte(err.Error())
			CreateJsonResponse(&w, false, &error_string)
			return
		}

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			var error_string = []byte(err.Error())
			CreateJsonResponse(&w, false, &error_string)
			return
		}

		CreateJsonResponse(&w, true, &bytes)

		return
	}

	data, err := pastelize.GeneratePalette(image)
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		log.Println("Failed to write color.json cache file. occurred on " + path)
	}

	CreateJsonResponse(&w, true, &bytes)
}
