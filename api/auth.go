package api

import (
	"net/http"
	"os"

	"github.com/CNAM-Projet-S3/api_storybuilder/services"
	"github.com/zmb3/spotify"
)

func IsValidClientID(clientID string) bool {
	return clientID == os.Getenv("API_CLIENT_ID")
}

func LoginSpotify(prv *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := services.TokenFromFile(); err == nil {
			w.Write([]byte("Spotify is already connected"))
			return
		}

		http.Redirect(w, r, prv.Auth.AuthURL(prv.State), 200)
	}
}

func CompleteSpotifyAuth(prv *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := services.TokenFromFile(); err == nil {
			w.Write([]byte("Spotify is already connected"))
			return
		}

		token, err := prv.Auth.Token(prv.State, r)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if st := r.FormValue("state"); st != prv.State {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		client := spotify.Authenticator{}.NewClient(token)
		tok, err := client.Token()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = services.TokenToFile(tok)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Write([]byte("You can close this window"))

	}
}
