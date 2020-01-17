package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/CNAM-Projet-S3/api_storybuilder/services"
)

type HereRequest struct {
	ClientID  string
	Longitude string
	Latitude  string
}

type hereItem struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

type HereResponse struct {
	Results struct {
		Items []hereItem `json:"items"`
	} `json:"results"`
}

type HereApiResponse struct {
	Items []hereItem `json:"items"`
}

func HereEndpoint(prv *services.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var rc HereRequest

		err = json.Unmarshal(body, &rc)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(rc.ClientID) == 0 || len(rc.Longitude) == 0 || len(rc.Latitude) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !IsValidClientID(rc.ClientID) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		url := "https://places.cit.api.here.com/places/v1/discover/around?app_id=" + os.Getenv("HERE_CLIENT_ID") + "&app_code=" + os.Getenv("HERE_CLIENT_CODE") + "&at="
		rsp, err := http.Get(url + rc.Latitude + "," + rc.Longitude)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err = ioutil.ReadAll(rsp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var hereResp HereResponse

		err = json.Unmarshal(body, &hereResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(HereApiResponse{
			Items: hereResp.Results.Items,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(json)
	}
}
