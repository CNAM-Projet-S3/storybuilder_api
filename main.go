package main

import (
	"log"
	"net/http"
	"os"
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
	r.HandleFunc("/spotify", api.SpotifyEndpoint).Methods("POST")
	r.HandleFunc("/spotify/login", api.LoginSpotify)
	r.HandleFunc("/spotify/callback", api.CompleteSpotifyAuth)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + os.Getenv("API_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
