package handlers

import (
	"io"
	"net/http"
	"os"
)

func GetIndex(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("./index.html")
	if err != nil {
		http.Error(w, "Failed to load index", http.StatusInternalServerError)
	}
	defer file.Close()

	var r io.Reader = file
	if _, err := io.Copy(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
