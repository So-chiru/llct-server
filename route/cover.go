package route

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/so-chiru/llct-server/utils"

	"golang.org/x/image/draw"
)

func isFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func getCoverFilePath(group_name string, size int, id string) (bool, string, string) {
	var base_folder = "./datas/" + group_name + "/" + id + "/"
	var cover = base_folder + "cover.jpg"

	if isFileExists(cover) {
		if size == 0 {
			return true, cover, ""
		}

		var sized = base_folder + "_cache/cover." + fmt.Sprint(size) + ".jpg.rsz"
		return isFileExists(sized), cover, sized
	}

	return false, "", ""
}

func openImage(path string) (image.Image, error) {
	fl, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(fl)
	if err != nil {
		return nil, err
	}

	defer fl.Close()

	return img, nil
}

func resizeImage(origin_path string, path string, size int, w io.Writer) error {
	img, err := openImage(origin_path)
	if err != nil {
		return err
	}

	var x = img.Bounds().Max.X
	var y = img.Bounds().Max.Y

	var ratio_x float32 = float32(size) / float32(x)
	var ratio_y float32 = float32(size) / float32(y)

	x = int(float32(x) * ratio_x)
	y = int(float32(y) * ratio_y)

	dr := image.Rect(0, 0, x, y)
	var res image.Image = scaleTo(img, dr, draw.CatmullRom)

	var base = filepath.Dir(path) + "/"
	if !checkFileExists(base) {
		err := os.MkdirAll(base, 0777)
		if err != nil {
			return err
		}
	}

	writer, err := os.Create(path)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(w, writer)
	err = jpeg.Encode(mw, res, &jpeg.Options{
		Quality: 80,
	})

	defer writer.Close()

	if err != nil {
		return err
	}

	return nil
}

func scaleTo(src image.Image,
	rect image.Rectangle, scale draw.Scaler) image.Image {
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
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

	var size = r.URL.Query().Get("s")
	var size_number = 0

	if len(size) > 1 {
		sn, err := strconv.Atoi(size)

		if err != nil || sn >= 600 {
			sn = 0
		}

		if sn%25 != 0 {
			if sn%25 > 12 {
				sn = 25 * ((sn / 25) + 1)
			} else {
				sn = 25 * (sn / 25)
			}
		}

		size_number = sn
	}

	cache_exists, origin_path, sized_path := getCoverFilePath(group_name, size_number, fmt.Sprint(song_number))

	if len(origin_path) < 1 {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")

	fmt.Println(cache_exists, origin_path, sized_path)

	// 캐시 파일이 존재하는 경우
	if cache_exists && len(sized_path) > 1 {
		err := cacheFile(w, r, sized_path)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
		}

		return
	}

	// 원본 사이즈를 요청하는 경우
	if size_number == 0 {
		cacheFile(w, r, origin_path)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
		}

		return
	}

	err = resizeImage(origin_path, sized_path, size_number, w)

	if err != nil {
		log.Println(err)

		w.WriteHeader(404)
		return
	}
}
