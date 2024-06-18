package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sngbd/last-boxd-server/lib"
)

type Response struct {
	ImageBase64 string `json:"image"`
}

func LastBoxd(w http.ResponseWriter, r *http.Request) {
	var (
		col         int    = 3
		row         int    = 3
		qTitle      string = "on"
		qDirector   string = "on"
		qRating     string = "on"
		timeInt     int
		timeRange   time.Time
		imageBase64 string
		err         error
	)

	log.Printf("%s", r.URL.String())

	vars := mux.Vars(r)
	q := r.URL.Query()

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

	if qRow == "" && qCol == "" {
		col, row = 0, 0
		timeInt, err = strconv.Atoi(q.Get("time"))
		if err != nil {
			log.Fatal(err)
		}

		loc, _ := time.LoadLocation(time.UTC.String())
		currentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, loc)
		switch timeInt {
		case 1:
			timeRange = currentTime.AddDate(0, 0, -7)
		case 2:
			timeRange = currentTime.AddDate(0, -1, 0)
		case 3:
			timeRange = currentTime.AddDate(0, -3, 0)
		case 4:
			timeRange = currentTime.AddDate(-1, 0, 0)
		}
	}

	if col == 0 && row == 0 {
		imageBase64 = lib.GetLastBoxdTime(username, timeRange, qTitle, qDirector, qRating)
	} else {
		imageBase64 = lib.GetLastBoxd(username, col, row, qTitle, qDirector, qRating)
	}

	data := Response{ImageBase64: imageBase64}

	json, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(json)
}
