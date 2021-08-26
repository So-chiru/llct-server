package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/so-chiru/llct-server/metaify/metautils"
	"github.com/so-chiru/llct-server/structs"
	"github.com/so-chiru/llct-server/utils"
	"github.com/tcolgate/mp3"
)

func divider() {
	writeLog("------------------------------------------------")
}

func printHelp() {
	println("Usage: metaify mode (id) [--display] [--skip-input] [--override]")
	println()
	println("\t--display: Skip file save and print serialized JSON data to stdout.\n")
	println("\t--skip-input: Request to not ask the input on the program.\n")
	println("\t--override: Allow program to override things automatically if program needs to save something.\n")
	println("  metadata:")
	println("\tGenerate metadata based on audio file and write it to lists.json file.")
	println("  new:")
	println("\tCreate new song and write it to lists.json file.")
	println("  cover:")
	println("\tDownload cover from iTunes server.")
	println("  normalize:")
	println("\tNormalize audio file.")
	os.Exit(0)
}

func getMetadata(dir string) (tag.Metadata, error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	meta, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func getDuration(dir string) (float64, error) {
	file, err := os.Open(dir)
	if err != nil {
		return 0, err
	}

	skipped := 0
	duration := 0.0
	d := mp3.NewDecoder(file)
	var f mp3.Frame
	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				return duration, nil
			}

			return 0, err
		}

		duration += f.Duration().Seconds()
	}
}

func timeReleasedString(s string) int64 {
	var layout = "20060102"

	if strings.Contains(s, "년") {
		layout = "2006년 01월 02일"
	} else if strings.Contains(s, "-") {
		layout = "2006-01-02"
	}

	t, err := time.Parse(layout, s)
	if err != nil {
		log.Panicln(err)
	}

	return t.UnixNano() / int64(time.Millisecond) / 1000
}

var displayMode = false
var doOverride = false
var noInput = false

func writeLog(a ...interface{}) {
	if displayMode {
		return
	}

	fmt.Fprintln(os.Stdout, a...)
}

func getId() (string, int, int) {
	var id = os.Args[2]

	if len(id) > 5 || len(id) <= 1 {
		panic("Invalid ID length.")
	}

	if _, err := strconv.Atoi(id); err != nil {
		errString := fmt.Sprintf("Given id %q is not type of int.", id)
		panic(errString)
	}

	groupId, err := strconv.Atoi(id[:1])
	if err != nil {
		panic("Invalid group ID.")
	}
	songId, err := strconv.Atoi(id[1:])
	if err != nil {
		panic("Invalid song ID.")
	}

	return id, groupId, songId
}

func askContinue() {
	var result = ""
	_, err := fmt.Scanln(&result)
	if err != nil {
		if err.Error() == "unexpected newline" {
			result = "Y"
		} else {
			panic(err)
		}
	}

	if result != "Y" && result != "y" {
		writeLog("[-] Aborts!")
		os.Exit(0)
	}
}

func generateMetadata() {
	if len(os.Args) < 3 {
		fmt.Println("ID value is not defined.")
		os.Exit(1)
	}

	id, groupId, songId := getId()

	var lists = utils.GetListFile()

	var group = lists.Groups[groupId]
	var song = lists.Songs[groupId][songId-1]

	writeLog("Group: " + group.Name + " (" + group.Id + ")")
	writeLog("Song: " + song.Title + "")

	divider()

	if song.Metadata != nil {
		writeLog("[-] Album: ", song.Metadata.Album)
		writeLog("[-] BPM: ", song.Metadata.BPM)
		writeLog("[-] Composer: ", song.Metadata.Composer)
		writeLog("[-] Length: ", song.Metadata.Length)
		writeLog("[-] Released: ", song.Metadata.Released)

		if !doOverride {
			fmt.Print("[+] Album metadata is already exists. Do you want to override? [Y/n] ")
			askContinue()
		}

	} else {
		writeLog("[-] Empty metadata.")
	}

	divider()

	var audioPath = "./datas/" + group.Id + "/" + id[1:] + "/audio.mp3"
	meta, err := getMetadata(audioPath)
	if err != nil {
		log.Fatalln(err)
	}

	duration, err := getDuration(audioPath)
	if err != nil {
		log.Fatalln(err)
	}

	var purifiedDuration = int64(duration)

	writeLog("[+] Audio duration:", purifiedDuration)

	var album = meta.Album()
	var composer = meta.Composer()

	writeLog("[+] New album:", album)
	writeLog("[+] New composer:", composer)

	var newMetadata = &structs.LLCTSongMetadata{
		Album:  album,
		Length: purifiedDuration,
	}

	if !noInput {
		var dateScan = metautils.AskNewInput("Released date (YYYYmmdd)")
		var bpmScan = metautils.AskNewInput("BPM")

		if len(dateScan) > 0 {
			var date = timeReleasedString(dateScan)
			newMetadata.Released = date
		}

		if len(bpmScan) > 0 {
			var bpm = 0
			if len(bpmScan) != 0 {
				bpm, err = strconv.Atoi(bpmScan)
				if err != nil {
					log.Fatalln(err)
				}
			}

			newMetadata.BPM = int64(bpm)
		}
	}

	if len(composer) != 0 {
		newMetadata.Composer = strings.Split(composer, ",")

		for i := 0; i < len(newMetadata.Composer); i++ {
			if len(newMetadata.Composer[i]) < 1 {
				continue
			}

			newMetadata.Composer[i] = strings.Trim(newMetadata.Composer[i], " ")
		}
	}

	lists.Songs[groupId][songId-1].Metadata = newMetadata

	var buf = utils.ToListFile(lists, displayMode)

	if displayMode {
		fmt.Println(buf.String())
	}
}

func createNewSong() {
	id, groupId, songId := getId()

	var lists = utils.GetListFile()

	var group = lists.Groups[groupId]

	if len(lists.Songs[groupId]) > songId-1 {
		writeLog("Already exists. (" + lists.Songs[groupId][songId-1].Title + "). ID " + fmt.Sprint(groupId) + fmt.Sprint(len(lists.Songs[groupId])+1) + " is available for use.")
		os.Exit(0)
	}

	var song = &structs.LLCTSongs{}

	var title = metautils.AskNewInput("Title")

	song.Title = title

	if len(title) == 0 {
		writeLog("You must define a title.")
		os.Exit(1)
	}

	var titleKo = metautils.AskNewInput("TitleKo")

	if len(titleKo) != 0 {
		song.TitleKo = &titleKo
	}

	var artistHelpText = "Artist (\n"
	for i := 0; i < len(group.Artists); i++ {
		artistHelpText += fmt.Sprint(i) + ": " + group.Artists[i] + "\n"
	}
	artistHelpText += "), default 0"

	var artistText = metautils.AskNewInput(artistHelpText)

	if len(artistText) == 0 {
		artistText = "0"
	}

	artist, err := strconv.Atoi(artistText)
	if err != nil {
		panic(err)
	}

	if artist >= len(group.Artists) {
		writeLog("Invalid artist value. should less or same than " + fmt.Sprint(len(group.Artists)-1))
		os.Exit(1)
	}

	song.Artist = artist

	lists.Songs[groupId] = append(lists.Songs[groupId], *song)

	utils.ToListFile(lists, displayMode)

	var audioPath = "./datas/" + group.Id + "/" + id[1:] + "/audio.mp3"

	if _, err := os.Open(audioPath); os.IsExist(err) {
		generateMetadata()
	}
}

func coverDownload() {
	id, groupId, songId := getId()

	var lists = utils.GetListFile()

	var group = lists.Groups[groupId]
	var song = lists.Songs[groupId][songId-1]

	var title = song.Title

	var coverPath = "./datas/" + group.Id + "/" + id[1:] + "/cover.jpg"

	if _, err := os.Stat(coverPath); err == nil {
		fmt.Print("[-] Cover for ID " + id + " (" + title + ") is already exists. Do you want to override? [Y/n] ")
		askContinue()
	}

	url, err := metautils.SearchiTunesCover(title)
	if err != nil {
		panic(err)
	}

	writeLog("[+] Downloading", url)

	r, err := metautils.DownloadCover(url)
	if err != nil {
		panic(err)
	}

	writeLog("[+] Saving", url)
	writeLog("[+] Destination:", coverPath)

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(coverPath, bytes, 0655)
	if err != nil {
		panic(err)
	}

	writeLog("[+] File saved.")
}

func normalizeAudio() {
	id, groupId, songId := getId()

	var lists = utils.GetListFile()

	var group = lists.Groups[groupId]
	var song = lists.Songs[groupId][songId-1]

	var audioPath = "./datas/" + group.Id + "/" + id[1:] + "/audio.mp3"

	writeLog("[+] Normalizing " + id + " (" + song.Title + ")")

	cmd := exec.Command("ffmpeg-normalize", audioPath, "-c:a", "mp3", "-o", audioPath, "-ar 44100", "-b:a 192k", "-f", "-nt", "rms", "-t", "-15")
	so, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	eo, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, so)
	// io.Copy(os.Stderr, eo)

	b, err := ioutil.ReadAll(eo)
	if err != nil {
		panic(err)
	}

	if len(b) != 0 {
		fmt.Print(string(b))
		writeLog("[+] Normalize failed.")
	} else {
		writeLog("[+] Normalize done.")
	}
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
	}

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--display" {
			displayMode = true
		} else if os.Args[i] == "--skip-input" {
			noInput = true
		} else if os.Args[i] == "--override" {
			doOverride = true
		} else if strings.Index(os.Args[i], "--") == 0 {
			println("Unknown argument:", os.Args[i]+"\n")
			printHelp()
		}

		if os.Args[i] == "--help" {
			printHelp()
		}
	}

	if os.Args[1] == "metadata" {
		generateMetadata()
	} else if os.Args[1] == "new" {
		createNewSong()
	} else if os.Args[1] == "cover" {
		coverDownload()
	} else if os.Args[1] == "normalize" {
		normalizeAudio()
	} else {
		printHelp()
	}
}
