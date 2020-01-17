package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type Provider struct {
	Auth  spotify.Authenticator
	State string

	Token *oauth2.Token
}

func TokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open("spotify.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func TokenToFile(token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", "spotify.json")
	f, err := os.OpenFile("spotify.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}
