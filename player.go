package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name  string 
	IsDir bool
	Mode  os.FileMode
}

const (
	filePrefix = "/music/"
	root = "./music"
)



func main() {
	http.HandleFunc("/", playerMainFrame)
	http.HandleFunc(filePrefix, File)
	http.ListenAndServe(":8080", nil)
}

func playerMainFrame(w http.ResponseWriter, r *http.Request) {
	if (r.URL.Path == "/") {
		http.ServeFile(w, r, "./player.html")///// + r.URL.Path)
	} else {
		http.ServeFile(w, r, "." + r.URL.Path)
	}
}

func File(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(root, r.URL.Path[len(filePrefix):])
	stat, err := os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if stat.IsDir() {
		serveDir(w, r, path)
		return
	}
	http.ServeFile(w, r, path)
}

func serveDir(w http.ResponseWriter, r *http.Request, path string) {
	defer func() {
		if err, ok := recover().(error); ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()
	file, err := os.Open(path)
	defer file.Close();
	if err != nil {
		panic(err)
	}

	files, err := file.Readdir(-1)
	if err != nil {
		panic(err)
	}

	fileinfos := make([]FileInfo, len(files), len(files))

	for i := range files {
		fileinfos[i].Name  = files[i].Name()
		fileinfos[i].IsDir = files[i].IsDir()
		fileinfos[i].Mode  = files[i].Mode()
	}

	j := json.NewEncoder(w)

	if err := j.Encode(&fileinfos); err != nil {
		panic(err)
	}
}
