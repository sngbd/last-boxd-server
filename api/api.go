package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sngbd/last-boxd/lib"
)

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

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<img src="data:image/jpeg;base64,%s">`, imageBase64)
}
