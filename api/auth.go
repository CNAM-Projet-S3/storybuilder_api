package api

import "os"

func IsValidClientID(clientID string) bool {
	return clientID == os.Getenv("API_CLIENT_ID")
}
