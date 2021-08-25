package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/so-chiru/llct-server/structs"
)

// 주어진 int 값의 ID가 lists 파일에 존재하는지 확인하여
// exists: bool, id: string 형식의 두 값으로 반환합니다.
func GetGroupID(i int) (bool, string) {
	var data = GetListFile()

	if i < 0 || i >= len(data.Groups) {
		return false, ""
	}

	return true, data.Groups[i].Id
}

func GetListFile() structs.LLCTLists {
	bytes, err := os.ReadFile("./datas/lists.json")

	if err != nil {
		log.Panicln("Couldn't read lists.json file.")
	}

	var data (structs.LLCTLists)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Panicln("lists.json file is not valid JSON file.", err)
	}

	return data
}

func ToListFile(list structs.LLCTLists, noSave bool) bytes.Buffer {
	data, err := json.Marshal(list)
	if err != nil {
		log.Panicln("Given list data is not valid structs.LLCTLists struct.", err)
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, data, "", "  ")
	if err != nil {
		log.Panicln("Failed to prettify.", err)
	}

	if !noSave {
		err = ioutil.WriteFile("./datas/lists.json", pretty.Bytes(), 0755)
		if err != nil {
			log.Panicln("Failed to write lists.json file", err)
		}
	}

	return pretty
}
