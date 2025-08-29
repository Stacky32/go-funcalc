package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
)

func GetIndex(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("./index.html")
	if err != nil {
		wd, _ := os.Getwd()
		log.Println("CWD", wd)
		log.Printf("Failed to open index.html: %#v", err)
		http.Error(w, "Failed to load index", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	var r io.Reader = file
	if _, err := io.Copy(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
