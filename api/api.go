package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sngbd/last-boxd/lib"
)

type Response struct {
	ImageBase64 string `json:"image"`
}

func LastBoxd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	q := r.URL.Query()

	var grid int = 3
	var text string = "on"
	var err error

	username := vars["username"]
	qGrid := q.Get("grid")
	text = q.Get("text")

	if qGrid != "" {
		grid, err = strconv.Atoi(qGrid)
		if err != nil {
			log.Fatal(err)
		}
	}

	imageBase64 := lib.GetLastBoxd(username, grid, text)

	data := Response{ImageBase64: imageBase64}

	json, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(json)
}
