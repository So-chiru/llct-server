package route

import (
	"log"
	"net/http"
	"strconv"

	"io/ioutil"
	"os"
	"time"
)

var cached_lists *[]byte
var cached_lists_date time.Time

func listsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// 마지막 디스크 읽기 후 5초가 지나지 않았으면 디스크에서 불러오지 않고 메모리에서 읽음.
	if cached_lists != nil && time.Now().Before(cached_lists_date) {
		CreateJsonResponse(&w, true, cached_lists)

		return
	}

	f, err := os.Open("./datas/lists.json")
	if err != nil {
		var error_string []byte = []byte("리스트 파일을 디스크에서 로딩하는 중 오류가 발생했습니다.")
		CreateJsonResponse(&w, false, &error_string)

		log.Println(err)
		return
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		var error_string []byte = []byte("버퍼를 읽을 수 없습니다.")
		CreateJsonResponse(&w, false, &error_string)

		log.Println(err)
		return
	}

	// 메모리에 캐싱
	cached_lists = &bytes

	t, err := strconv.Atoi(os.Getenv("CACHE_DURATION"))
	if err != nil {
		log.Println("The environmental variable 'CACHE_DURATION' is invalid or not defined.")
		t = 5
	}

	cached_lists_date = time.Now().Add(time.Second * time.Duration(t))

	CreateJsonResponse(&w, true, cached_lists)
}
