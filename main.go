package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sngbd/last-boxd/api"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	port := fmt.Sprint(viper.Get("PORT"))

	router := mux.NewRouter()

	router.HandleFunc("/{username}", api.LastBoxd).Methods("GET")

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
