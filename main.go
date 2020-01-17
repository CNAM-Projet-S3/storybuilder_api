package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/oxodao/api_storybuilder/api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/here", api.HereEndpoint).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:1242",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
