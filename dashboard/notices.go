package dashboard

import (
	"encoding/json"
	"os"
)

func openNoticesFile() (*[]byte, error) {
	bytes, err := os.ReadFile("./datas/notices.json")

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &bytes, nil
}

func GetNoticesData() (*[]Notices, error) {
	bytes, err := openNoticesFile()
	if err != nil {
		return nil, err
	}

	var noti []Notices
	err = json.Unmarshal(*bytes, &noti)
	if err != nil {
		return nil, err
	}

	return &noti, nil
}
