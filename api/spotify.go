package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/CNAM-Projet-S3/api_storybuilder/services"
	"github.com/zmb3/spotify"
)

type SpotifyRequest struct {
	ClientID string

	Artist string
	Album  string
	Title  string

	ID string
}

type SpotifyResponse struct {
	Artist string
	Title  string
	ID     string
	Cover  string
}

func SpotifyEndpoint(prv *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

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

		//token, err := prv.Auth.Exchange(acTk)
		token, err := services.TokenFromFile()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Can't load token"))

			return
		}

		client := spotify.Authenticator{}.NewClient(token)

		// search for playlists and albums containing "holiday"
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// If there is an ID, don't bother searching it and directly responding whats needed
		if len(rc.ID) != 0 {

			return
		} else {
			var searchtype spotify.SearchType = spotify.SearchTypeArtist
			query := rc.Artist + " "

			// If there is an album, search for it, fallback on the song
			if len(rc.Album) != 0 {
				searchtype = searchtype | spotify.SearchTypeAlbum
				query = query + rc.Album
			} else {
				searchtype = searchtype | spotify.SearchTypeTrack
				query = query + rc.Title
			}

			results, err := client.Search(query, searchtype)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			if results.Tracks != nil {
				trc := results.Tracks.Tracks[0]
				rsp := SpotifyResponse{
					Artist: trc.Artists[0].Name,
					Title:  trc.Name,
					ID:     trc.ID.String(),
					Cover:  trc.Album.Images[0].URL,
				}

				json, err := json.Marshal(rsp)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Write(json)
				return
			}

			if results.Albums != nil {
				alb := results.Albums.Albums[0]
				rsp := SpotifyResponse{
					Artist: alb.Artists[0].Name,
					Title:  alb.Name,
					ID:     alb.ID.String(),
					Cover:  alb.Images[0].URL,
				}

				json, err := json.Marshal(rsp)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Write(json)
				return
			}

		}

	}
}
