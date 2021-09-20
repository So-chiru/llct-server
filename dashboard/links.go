package dashboard

import (
	"encoding/json"
	"os"
)

func openLinksFile() (*[]byte, error) {
	bytes, err := os.ReadFile("./datas/links.json")

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &bytes, nil
}

func GetLinksData() (*[]LinkData, error) {
	bytes, err := openLinksFile()
	if err != nil {
		return nil, err
	}

	var links []LinkData
	err = json.Unmarshal(*bytes, &links)
	if err != nil {
		return nil, err
	}

	return &links, nil
}
