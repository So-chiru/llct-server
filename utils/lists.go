package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/so-chiru/llct-server/structs"
)

// 주어진 int 값의 ID가 lists 파일에 존재하는지 확인하여
// exists: bool, id: string 형식의 두 값으로 반환합니다.
func GetGroupID(i int) (bool, string) {
	bytes, err := os.ReadFile("./datas/lists.json")

	if err != nil {
		log.Panicln("Couldn't read lists.json file.")
	}

	var data (structs.LLCTLists)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Panicln("lists.json file is not valid JSON file.", err)
	}

	if i < 0 || i >= len(data.Groups) {
		return false, ""
	}

	return true, data.Groups[i].Id
}
