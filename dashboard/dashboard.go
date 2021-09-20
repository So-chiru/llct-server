package dashboard

import (
	"fmt"
	"time"

	"github.com/so-chiru/llct-server/characters"
)

type Notices struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type BirthdayComponent struct {
	Name   string   `json:"name"`
	Color  string   `json:"color"`
	Date   string   `json:"date"`
	Musics []string `json:"musics"`
}

type MusicComponent struct {
	ID              string  `json:"id"`
	RecommendReason *string `json:"recommendReason,omitempty"`
}

type LinkComponent struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Link        string  `json:"link"`
}

type LinkData struct {
	Title string          `json:"title"`
	Items []LinkComponent `json:"items"`
}

type LiveComponent struct {
	Title      *string  `json:"title,omitempty"`
	Image      string   `json:"image"`
	Start      string   `json:"start"`
	End        string   `json:"end"`
	URL        *string  `json:"url,omitempty"`
	Location   *string  `json:"location,omitempty"`
	Characters []string `json:"characters"`
}

type Dashboard struct {
	Type      string             `json:"type"`
	Title     *string            `json:"title,omitempty"`
	Birthday  *BirthdayComponent `json:"birthday,omitempty"`
	Music     *MusicComponent    `json:"music,omitempty"`
	Link      *LinkComponent     `json:"link,omitempty"`
	Live      *LiveComponent     `json:"live,omitempty"`
	CustomSet *[]interface{}     `json:"customset,omitempty"`
	MusicSet  *[]MusicComponent  `json:"musicset,omitempty"`
	LinkSet   *[]LinkComponent   `json:"linkset,omitempty"`
}

func generateBirthday() *[]Dashboard {
	birthdays, err := GetBirthdayData()
	if err != nil {
		return nil
	}

	if birthdays == nil {
		return nil
	}

	var result []Dashboard

	for _, v := range *birthdays {
		t, err := time.Parse("2006-01-02", v.Start)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if t.Day() == time.Now().Day() && t.Month() == time.Now().Month() {
			c, err := characters.GetCharacterData(v.Character)

			if c == nil || err != nil {
				continue
			}

			d := t.Format(time.RFC3339)

			var birth = &BirthdayComponent{
				Name:   v.Character,
				Color:  c.Color,
				Date:   d,
				Musics: c.Musics,
			}

			var data = Dashboard{
				Type:     "birthday",
				Birthday: birth,
			}
			result = append(result, data)
		}

	}

	return &result
}

func generateLive() *[]Dashboard {
	lives, err := GetLiveData()
	if err != nil {
		return nil
	}

	if lives == nil {
		return nil
	}

	var result []Dashboard

	for _, v := range *lives {
		_, err := time.Parse(time.RFC3339, v.Start)
		if err != nil {
			fmt.Println(err)
			continue
		}

		e, err := time.Parse(time.RFC3339, v.End)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if time.Now().Before(e) || time.Now().Day() == e.Day() {
			var title = v.Name

			var live = LiveComponent{
				Title:      &title,
				Image:      "unsupported",
				Start:      v.Start,
				End:        v.End,
				URL:        v.URL,
				Location:   v.Location,
				Characters: v.Characters,
			}

			var data = Dashboard{
				Type:  "live",
				Title: &title,
				Live:  &live,
			}
			result = append(result, data)
		}

	}

	return &result
}

func generateLinks() *[]Dashboard {
	links, err := GetLinksData()
	if err != nil {
		return nil
	}

	if links == nil {
		return nil
	}

	var result []Dashboard

	for _, link := range *links {
		var data = Dashboard{
			Type:    "linkset",
			Title:   &link.Title,
			LinkSet: &link.Items,
		}
		result = append(result, data)
	}

	return &result
}

func GetDashboards() *[]Dashboard {
	var result []Dashboard = make([]Dashboard, 0)

	birthday := generateBirthday()
	if birthday != nil {
		result = append(result, *birthday...)
	}

	randomSongs := generateRandomSongs(4)
	if birthday != nil {
		result = append(result, *randomSongs...)
	}

	links := generateLinks()
	if links != nil {
		result = append(result, *links...)
	}

	live := generateLive()
	if live != nil {
		result = append(result, *live...)
	}

	return &result
}

func GetNotices() *[]Notices {
	var result []Notices = make([]Notices, 0)

	notices, _ := GetNoticesData()
	if notices != nil {
		result = append(result, *notices...)
	}

	return &result
}
