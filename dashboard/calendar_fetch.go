package dashboard

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var userAgent = "LLCT/1.0"

type LLCalendarBirthdayResponse struct {
	Id        string `json:"id"`
	Character string `json:"character"`
	Start     string `json:"start"`
}

type LLCalendarLiveResponse struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Characters []string `json:"characters"`
	Start      string   `json:"start"`
	End        string   `json:"end"`
	URL        *string  `json:"url,omitempty"`
	Location   *string  `json:"location,omitempty"`
}

func fetchBirthdayCalendar() (*[]LLCalendarBirthdayResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://cal-api.llasfans.net/api/llct/birthday", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result *[]LLCalendarBirthdayResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("failed to fetch birthday calendar data. API returned nil")
	}

	return result, nil
}

func fetchLiveCalendar() (*[]LLCalendarLiveResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://cal-api.llasfans.net/api/llct/live", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result *[]LLCalendarLiveResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("failed to fetch birthday calendar data. API returned nil")
	}

	return result, nil
}
