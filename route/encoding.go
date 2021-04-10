package route

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/andybalholm/brotli"
)

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func cacheFile(w http.ResponseWriter, r *http.Request, path string) error {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
		w.Header().Set("Content-Encoding", "br")

		var cache_path string
		if strings.Contains(path, "_cache") {
			cache_path = filepath.Dir(path) + "/" + filepath.Base(path) + ".br"
		} else {
			cache_path = filepath.Dir(path) + "/_cache/" + filepath.Base(path) + ".br"
		}

		if checkFileExists(cache_path) {
			reader, err := os.Open(cache_path)
			if err != nil {
				return err
			}

			io.Copy(w, reader)

			defer reader.Close()

			return nil
		}

		// 파일 Reader
		reader, err := os.Open(path)
		if err != nil {
			return err
		}

		if !checkFileExists(filepath.Dir(cache_path)) {
			err := os.MkdirAll(filepath.Dir(cache_path), 0777)
			if err != nil {
				return err
			}
		}

		file, err := os.Create(cache_path)
		if err != nil {
			return err
		}

		// Writer를 2개로 나눔 : http.ResponseWriter, os.FileWriter
		writer := io.MultiWriter(w, file)
		br := brotli.NewWriterLevel(writer, brotli.BestCompression)
		io.Copy(br, reader)

		defer br.Close()

		return nil
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")

		var cache_path string
		if strings.Contains(path, "_cache") {
			cache_path = filepath.Dir(path) + "/" + filepath.Base(path) + ".gz"
		} else {
			cache_path = filepath.Dir(path) + "/_cache/" + filepath.Base(path) + ".gz"
		}

		if checkFileExists(cache_path) {
			reader, err := os.Open(cache_path)
			if err != nil {
				return err
			}

			io.Copy(w, reader)

			defer reader.Close()

			return nil
		}

		// 파일 Reader
		reader, err := os.Open(path)
		if err != nil {
			return err
		}

		if !checkFileExists(filepath.Dir(cache_path)) {
			err := os.MkdirAll(filepath.Dir(cache_path), 0777)
			if err != nil {
				return err
			}
		}

		file, err := os.Create(cache_path)
		if err != nil {
			return err
		}

		// Writer를 2개로 나눔 : http.ResponseWriter, os.FileWriter
		writer := io.MultiWriter(w, file)
		gz, err := gzip.NewWriterLevel(writer, gzip.BestCompression)
		if err != nil {
			return err
		}

		io.Copy(gz, reader)

		defer gz.Close()

		return nil
	}

	// 파일 Reader
	reader, err := os.Open(path)
	if err != nil {
		return err
	}

	io.Copy(w, reader)

	defer reader.Close()

	return nil
}
