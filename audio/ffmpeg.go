package audio

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func Compress(p string, codec string, w io.Writer, to string) (io.ReadSeeker, error) {
	// var cmds = []string{"-i", "-map", "0:a:0", "-b:a", fmt.Sprint(to) + "k,"}

	var cache_path = filepath.Dir(p) + "/_cache/" + filepath.Base(p) + "." + to

	var cmd *exec.Cmd

	if codec == "opus" {
		cmd = exec.Command("ffmpeg", "-i", "-", "-f", "opus", "-map_metadata", "-1", "-map", "0:a:0", "-b:a", to+"k", cache_path)
	} else {
		cmd = exec.Command("ffmpeg", "-i", "-", "-f", "mp3", "-map_metadata", "-1", "-map", "0:a:0", "-b:a", to+"k", cache_path)
	}

	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	return err
	// }

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	reader, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	stdin.Write(bytes)
	defer stdin.Close()

	reader, err = os.Open(cache_path)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
