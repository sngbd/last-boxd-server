package api

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func LastBoxd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	GetLastBoxd(username)

	// Open the image file
	f, err := os.Open("grid.jpg")
	if err != nil {
		http.Error(w, "file not found", 404)
		return
	}
	defer f.Close()

	// Set the content type to JPEG
	w.Header().Set("Content-Type", "image/jpeg")

	// Copy the image data to the response
	io.Copy(w, f)
}
