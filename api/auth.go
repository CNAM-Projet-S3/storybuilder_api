package api

import (
	"net/http"
	"os"

	"github.com/zmb3/spotify"
)

var (
	auth  spotify.Authenticator
	state = "toto123"
)

func IsValidClientID(clientID string) bool {
	return clientID == os.Getenv("API_CLIENT_ID")
}

func LoginSpotify(w http.ResponseWriter, r *http.Request) {
	acTk := os.Getenv("SPOTIFY_ACCESS_TOKEN")
	rfTk := os.Getenv("SPOTIFY_REFRESH_TOKEN")

	if len(acTk) == 0 || len(rfTk) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The server is already authentificated"))
		return
	}

	auth = spotify.NewAuthenticator(os.Getenv("HOST")+"/spotify/callback", spotify.ScopeUserReadPrivate)
	auth.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_KEY"))

	http.Redirect(w, r, auth.AuthURL(state), 200)
}

func CompleteSpotifyAuth(w http.ResponseWriter, r *http.Request) {
	acTk := os.Getenv("SPOTIFY_ACCESS_TOKEN")
	rfTk := os.Getenv("SPOTIFY_REFRESH_TOKEN")

	if len(acTk) == 0 || len(rfTk) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The server is already authentificated"))
		return
	}

	tok, err := auth.Token(state, r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if st := r.FormValue("state"); st != state {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte("Please store the token in the environment file: AccessToken = " + tok.AccessToken + " / RefreshToken = " + tok.RefreshToken))
}
