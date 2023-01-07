package main

import (
	"github.com/sngbd/last-boxd/api"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/{username}", api.LastBoxd).Methods("GET")

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println(err)
	}
}
