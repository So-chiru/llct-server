package dashboard

import (
	"log"
	"math"
	"os"
	"time"
)

func saveResponseCache(t string, r []byte) error {
	var file = "live_calendar"

	if t == "birthday" {
		file = "birthday_calendar"
	}

	err := os.WriteFile("./.cache/"+file, r, 0755)

	return err
}

func getResponseCache(t string) (*[]byte, error) {
	var file = "live_calendar"

	if t == "birthday" {
		file = "birthday_calendar"
	}

	bytes, err := os.ReadFile("./.cache/" + file)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &bytes, nil
}

func checkCacheAge(t string) (time.Duration, error) {
	var file = "live_calendar"

	if t == "birthday" {
		file = "birthday_calendar"
	}

	_, err := os.Stat("./.cache/")
	if os.IsNotExist(err) {
		err := os.Mkdir("./.cache/", 0755)

		if err != nil {
			panic(err)
		}

		log.Println("# created .cache dir")
	}

	stats, err := os.Stat("./.cache/" + file)

	if os.IsNotExist(err) {
		return math.MaxInt64, nil
	}

	if err != nil {
		return 0, err
	}

	return time.Since(stats.ModTime()), nil
}
