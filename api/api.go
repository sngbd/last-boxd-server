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
	username := vars["username"]
	varGrid := vars["grid"]

	var grid int = 3
	var err error

	if varGrid != "" {
		grid, err = strconv.Atoi(varGrid)
		if err != nil {
			log.Fatal(err)
		}
	}

	imageBase64 := lib.GetLastBoxd(username, grid)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<img src="data:image/jpeg;base64,%s">`, imageBase64)
}
