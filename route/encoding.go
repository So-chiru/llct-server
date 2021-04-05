package route

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/andybalholm/brotli"
)

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func cacheFile(w http.ResponseWriter, r *http.Request, path string) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
		w.Header().Set("Content-Encoding", "br")

		var cache_path = path + ".llct.br"

		if checkFileExists(cache_path) {
			reader, err := os.Open(cache_path)
			if err != nil {
				return
			}

			io.Copy(w, reader)

			return
		}

		// 파일 Reader
		reader, err := os.Open(path)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		file, err := os.Create(cache_path)
		if err != nil {
			return
		}

		// Writer를 2개로 나눔 : http.ResponseWriter, os.FileWriter
		writer := io.MultiWriter(w, file)
		br := brotli.NewWriterLevel(writer, brotli.BestCompression)
		io.Copy(br, reader)

		defer br.Close()

		return
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")

		var cache_path = path + ".llct.gz"

		if checkFileExists(cache_path) {
			reader, err := os.Open(cache_path)
			if err != nil {
				return
			}

			io.Copy(w, reader)

			return
		}

		// 파일 Reader
		reader, err := os.Open(path)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		file, err := os.Create(cache_path)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		// Writer를 2개로 나눔 : http.ResponseWriter, os.FileWriter
		writer := io.MultiWriter(w, file)
		gz, err := gzip.NewWriterLevel(writer, gzip.BestCompression)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		io.Copy(gz, reader)

		defer gz.Close()

		return
	}

	// 파일 Reader
	reader, err := os.Open(path)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	io.Copy(w, reader)
}
