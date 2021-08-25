package metautils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ItunesSongData struct {
	WrapperType            string  `json:"wrapperType"`
	CollectionType         string  `json:"collectionType"`
	ArtistId               float64 `json:"artistId"`
	CollectionId           float64 `json:"collectionId"`
	ArtistName             string  `json:"artistName"`
	CollectionName         string  `json:"collectionName"`
	CollectionCensoredName string  `json:"collectionCensoredName"`
	ArtistViewUrl          string  `json:"artistViewUrl"`
	CollectionViewUrl      string  `json:"collectionViewUrl"`
	ArtworkUrl60           string  `json:"artworkUrl60"`
	ArtworkUrl100          string  `json:"artworkUrl100"`
	CollectionPrice        float64 `json:"collectionPrice"`
	CollectionExplicitness string  `json:"collectionExplicitness"`
	TrackCount             float64 `json:"trackCount"`
	Copyright              string  `json:"copyright"`
	Country                string  `json:"county"`
	Currency               string  `json:"currency"`
	ReleaseDate            string  `json:"releaseDate"`
	PrimaryGenreName       string  `json:"primaryGenreName"`
}

type ITunesQueryData struct {
	ResultCount int              `json:"resultCount"`
	Results     []ItunesSongData `json:"results"`
}

func parseiTunesResponse(data []byte) (*ITunesQueryData, error) {
	var t = &ITunesQueryData{}

	err := json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func SearchiTunesCover(name string) (string, error) {
	var url = "https://itunes.apple.com/search?term=" + url.QueryEscape(name) + "&country=jp&entity=album&limit=25&_=1618203815335"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	data, err := parseiTunesResponse(bytes)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(data.Results); i++ {
		var item = data.Results[i]

		fmt.Println(i, item.CollectionName, item.ArtistName)
	}

	index := AskNewInput("Index to download (0)")
	if len(index) == 0 {
		index = "0"
	}

	indexInt, err := strconv.Atoi(index)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(data.Results[indexInt].ArtworkUrl100, "100x100", "600x600"), nil
}

func DownloadCover(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
