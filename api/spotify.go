package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type SpotifyRequest struct {
	ClientID string

	Artist string
	Album  string
	Title  string

	ID string
}

func SpotifyEndpoint(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rc SpotifyRequest

	err = json.Unmarshal(body, &rc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(rc.ClientID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !IsValidClientID(rc.ClientID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// If there is an ID, don't bother searching it and directly responding whats needed
	if len(rc.ID) != 0 {

		return
	}

	query := rc.Artist + " "

	// If there is an album, search for it, fallback on the song
	if len(rc.Album) != 0 {
		query = query + rc.Album
	} else {
		query = query + rc.Title
	}

}
