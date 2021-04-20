package route

import (
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/cenkalti/dominantcolor"
	"github.com/gorilla/mux"
	"github.com/so-chiru/llct-server/utils"
	"github.com/teacat/noire"
)

type SaturateOrder []noire.Color

func (a SaturateOrder) Len() int           { return len(a) }
func (a SaturateOrder) Less(i, j int) bool { return a[i].Saturation() < a[j].Saturation() }
func (a SaturateOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func colorHandler(w http.ResponseWriter, r *http.Request) {
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

	var data = dominantcolor.FindN(image, 5)

	colors := []noire.Color{
		noire.NewRGB(float64(data[0].R), float64(data[0].G), float64(data[0].B)),
		noire.NewRGB(float64(data[1].R), float64(data[1].G), float64(data[1].B)),
		noire.NewRGB(float64(data[2].R), float64(data[2].G), float64(data[2].B)),
		noire.NewRGB(float64(data[3].R), float64(data[3].G), float64(data[3].B)),
		noire.NewRGB(float64(data[4].R), float64(data[4].G), float64(data[4].B)),
	}

	sort.Sort(SaturateOrder(colors))

	var main = colors[2]
	var whiteTextColor = colors[2].Foreground()

	var mainDark = colors[4].Shade(0.4)
	var darkTextColor = colors[4].Foreground()

	if colors[2].IsDark() && colors[4].IsLight() {
		main = colors[4]
		whiteTextColor = colors[4].Foreground()
		mainDark = colors[2].Shade(0.4)
		darkTextColor = colors[2].Foreground()
	}

	var obj = GoColorStruct{
		Main:     "#" + main.Hex(),
		Sub:      "#" + main.Darken(0.1).Hex(),
		Text:     "#" + whiteTextColor.Hex(),
		MainDark: "#" + mainDark.Brighten(0.1).Hex(),
		SubDark:  "#" + mainDark.Hex(),
		TextDark: "#" + darkTextColor.Hex(),
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		var error_string = []byte(err.Error())
		CreateJsonResponse(&w, false, &error_string)
		return
	}

	CreateJsonResponse(&w, true, &bytes)
}
