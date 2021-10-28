package route

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/so-chiru/llct-server/audio"
	"github.com/so-chiru/llct-server/utils"
)

func getAudioFilePath(group_name string, id string) string {
	var base_folder = "./datas/" + group_name + "/" + id + "/"
	var cover = base_folder + "audio.mp3"

	if isFileExists(cover) {
		return cover
	}

	return ""
}

func audioHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var SERVE_AUDIO = os.Getenv("LLCT_SERVE_AUDIO")

	if SERVE_AUDIO != "1" && SERVE_AUDIO != "true" {
		var error_string = []byte("더 이상 음원을 제공하지 않습니다.")
		CreateJsonResponse(&w, false, &error_string)

		return
	}

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

	audio_path := getAudioFilePath(group_name, fmt.Sprint(song_number))

	if len(audio_path) < 1 {
		w.WriteHeader(404)

		return
	}

	stat, err := os.Stat(audio_path)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var bitrate = r.URL.Query().Get("b")

	if bitrate != "128" && bitrate != "64" {
		reader, err := os.Open(audio_path)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		if bitrate == "64" {
			w.Header().Set("Content-Type", "audio/opus")
		} else {
			w.Header().Set("Content-Type", "audio/mp3")
		}

		http.ServeContent(w, r, "llct.audio", stat.ModTime(), reader)

		return
	}

	var cache_path = filepath.Dir(audio_path) + "/_cache/" + filepath.Base(audio_path) + "." + bitrate
	if isFileExists(cache_path) {
		reader, err := os.Open(cache_path)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		http.ServeContent(w, r, "llct.audio", stat.ModTime(), reader)

		return
	}

	var codec string
	if bitrate == "64" {
		codec = "opus"
	} else {
		codec = "mp3"
	}

	if codec == "opus" {
		w.Header().Set("Content-Type", "audio/ogg")
	} else {
		w.Header().Set("Content-Type", "audio/mp3")
	}

	reader, err := audio.Compress(audio_path, codec, w, bitrate)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}

	http.ServeContent(w, r, "llct.audio", stat.ModTime(), reader)
}
