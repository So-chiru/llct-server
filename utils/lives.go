package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/so-chiru/llct-server/structs"
)

func GetLiveFile() map[string]structs.LLCTLive {
	bytes, err := os.ReadFile("./datas/lives.json")

	if err != nil {
		log.Panicln("Couldn't read lists.json file.")
	}

	var data (map[string]structs.LLCTLive)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Panicln("lives.json file is not valid JSON file.", err)
	}

	return data
}

func AddLiveData(t string) {
	lives := GetLiveFile()

	idb := sha256.Sum256([]byte(t))
	id := hex.EncodeToString(idb[:])

	_, exists := lives[id]
	if exists {
		return
	}

	var emptyPlaylist []structs.LLCTPlaylist = make([]structs.LLCTPlaylist, 0)

	for i := 0; i < 2; i++ {
		var template = "predict"

		if i == 1 {
			template = "actual"
		}

		var emptyMusics []string = make([]string, 0)

		var playlist = structs.LLCTPlaylist{
			Template: template,
			Musics:   emptyMusics,
		}

		emptyPlaylist = append(emptyPlaylist, playlist)
	}

	lives[id] = structs.LLCTLive{
		Title:     &t,
		Playlists: emptyPlaylist,
	}

	ToLiveFile(lives, false)
}

func GetLiveData(t string) *structs.LLCTLive {
	lives := GetLiveFile()

	idb := sha256.Sum256([]byte(t))
	id := hex.EncodeToString(idb[:])

	data, exists := lives[id]
	if !exists {
		return nil
	}

	return &data
}

func ToLiveFile(list map[string]structs.LLCTLive, noSave bool) bytes.Buffer {
	data, err := json.Marshal(list)
	if err != nil {
		log.Panicln("Given list data is not valid structs.LLCTLive struct.", err)
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, data, "", "  ")
	if err != nil {
		log.Panicln("Failed to prettify.", err)
	}

	if !noSave {
		err = ioutil.WriteFile("./datas/lives.json", pretty.Bytes(), 0755)
		if err != nil {
			log.Panicln("Failed to write live.json file", err)
		}
	}

	return pretty
}
