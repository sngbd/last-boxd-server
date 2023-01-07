package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func LastBoxd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	imageBase64 := GetLastBoxd(username)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<img src="data:image/jpeg;base64,%s">`, imageBase64)
}
