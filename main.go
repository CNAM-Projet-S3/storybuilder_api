package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CNAM-Projet-S3/api_storybuilder/api"
	"github.com/CNAM-Projet-S3/api_storybuilder/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	provider := services.Provider{
		Auth:  spotify.NewAuthenticator(os.Getenv("HOST")+"/spotify/callback", spotify.ScopeUserReadPrivate),
		State: "toto123",
	}

	provider.Auth.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_KEY"))

	r := mux.NewRouter()

	r.HandleFunc("/here", api.HereEndpoint(&provider)).Methods("POST")
	r.HandleFunc("/spotify", api.SpotifyEndpoint(&provider)).Methods("POST")
	r.HandleFunc("/spotify/login", api.LoginSpotify(&provider))
	r.HandleFunc("/spotify/callback", api.CompleteSpotifyAuth(&provider))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + os.Getenv("API_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
