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

	var col int = 3
	var row int = 3
	var qTitle string = "on"
	var qDirector string = "on"
	var qRating string = "on"
	var err error

	username := vars["username"]
	qCol := q.Get("col")
	qRow := q.Get("row")
	qTitle = q.Get("title")
	qDirector = q.Get("director")
	qRating = q.Get("rating")

	if qCol != "" {
		col, err = strconv.Atoi(qCol)
		if err != nil {
			log.Fatal(err)
		}
	}

	if qRow != "" {
		row, err = strconv.Atoi(qRow)
		if err != nil {
			log.Fatal(err)
		}
	}

	imageBase64 := lib.GetLastBoxd(username, col, row, qTitle, qDirector, qRating)

	data := Response{ImageBase64: imageBase64}

	json, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(json)
}
