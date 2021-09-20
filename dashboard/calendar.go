package dashboard

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/so-chiru/llct-server/utils"
)

func SyncBirthdayCalendar() {
	run := func(t *time.Time) {
		log.Println("# sync birthday calendar...")

		sinceCache, err := checkCacheAge("birthday")
		if err != nil {
			fmt.Println("# error occurred while checking cache age of birthday calendar", err)
			return
		}

		if sinceCache < time.Hour*24 {
			log.Println("# skipping fetching birthday calendar since age of cache is younger than 1 day")
			return
		}

		data, err := fetchBirthdayCalendar()
		if err != nil {
			fmt.Println("# failed to fetch birthday calendar data", err)
			return
		}

		bytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println("# failed to marshal birthday calendar data", err)
			return
		}

		err = saveResponseCache("birthday", bytes)
		if err != nil {
			fmt.Println("# failed to save birthday calendar data to cache", err)
			return
		}
	}

	ticker := time.NewTicker(time.Hour * 24 * 7)

	go func() {
		for t := range ticker.C {
			run(&t)
		}
	}()

	run(nil)
}

func SyncLiveCalendar() {
	run := func(t *time.Time) {
		log.Println("# sync live calendar...")

		sinceCache, err := checkCacheAge("live")
		if err != nil {
			fmt.Println("# error occurred while checking cache age of live calendar", err)
			return
		}

		if sinceCache < time.Hour*24 {
			log.Println("# skipping fetching live calendar since age of cache is younger than 1 day")
			return
		}

		data, err := fetchLiveCalendar()
		if err != nil {
			fmt.Println("# failed to fetch live calendar data", err)
			return
		}

		for _, v := range *data {
			utils.AddLiveData(v.Name)
		}

		bytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println("# failed to marshal live calendar data", err)
			return
		}

		err = saveResponseCache("live", bytes)
		if err != nil {
			fmt.Println("# failed to save live calendar data to cache", err)
			return
		}
	}

	ticker := time.NewTicker(time.Hour * 24 * 3)

	go func() {
		for t := range ticker.C {
			run(&t)
		}
	}()

	run(nil)
}

func GetBirthdayData() (*[]LLCalendarBirthdayResponse, error) {
	data, err := getResponseCache("birthday")
	if err != nil {
		return nil, err
	}

	var result *[]LLCalendarBirthdayResponse
	json.Unmarshal(*data, &result)

	return result, nil
}

func GetLiveData() (*[]LLCalendarLiveResponse, error) {
	data, err := getResponseCache("live")
	if err != nil {
		return nil, err
	}

	var result *[]LLCalendarLiveResponse
	json.Unmarshal(*data, &result)

	return result, nil
}

func GetCalendar() error {
	birthday, err := fetchBirthdayCalendar()
	if err != nil {
		return err
	}

	live, err := fetchLiveCalendar()
	if err != nil {
		return err
	}

	fmt.Println(birthday, live)

	return err
}
