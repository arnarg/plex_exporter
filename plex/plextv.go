package plex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PinRequest struct {
	Pin `json:"pin"`
}

type Pin struct {
	Id        int       `json:"id"`
	Code      string    `json:"code"`
	Expiry    time.Time `json:"expires_at"`
	Trusted   bool      `json:"trusted"`
	AuthToken string    `json:"auth_token"`
}

// GetPinRequest creates a PinRequest using the Plex API and returns it.
func GetPinRequest() (*PinRequest, error) {
	httpClient := &http.Client{Timeout: time.Second * 10}
	_, body, err := sendRequest("POST", "https://plex.tv/pins", headers, httpClient)
	if err != nil {
		return nil, err
	}

	pinRequest := &PinRequest{}

	err = json.Unmarshal(body, pinRequest)
	if err != nil {
		return nil, err
	}

	return pinRequest, nil
}

// GetTokenFromPinRequest takes in a PinRequest and checks if it has been authenticated.
// If it has been authenticated it returns the token.
// If it has not been authenticated it returns an empty string.
func GetTokenFromPinRequest(p *PinRequest) (string, error) {
	httpClient := &http.Client{Timeout: time.Second * 10}
	_, body, err := sendRequest("GET", fmt.Sprintf("https://plex.tv/pins/%d", p.Id), headers, httpClient)
	if err != nil {
		return "", err
	}

	pinRequest := PinRequest{}

	err = json.Unmarshal(body, &pinRequest)
	if err != nil {
		return "", err
	}

	if pinRequest.AuthToken == "" {
		return "", fmt.Errorf(ErrorPinNotAuthorized)
	}

	return pinRequest.AuthToken, nil
}
