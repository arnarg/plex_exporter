package plex

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/tagnard/plex_exporter/config"
	"github.com/tagnard/plex_exporter/plex/api"
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

func discoverServers(h map[string]string) ([]*Server, error) {
	httpClient := &http.Client{Timeout: time.Second * 10}
	// This endpoint only supports XML.
	// I want to specify the "Accept: application/xml" header
	// to make sure that if the endpoint does support JSON in
	// the future it won't break the application.
	eh := map[string]string{
		"Accept": "application/xml",
	}
	mergo.Merge(&eh, h)

	_, body, err := sendRequest("GET", "https://plex.tv/api/resources?includeHttps=1", eh, httpClient)
	if err != nil {
		return nil, err
	}

	deviceList := api.DeviceList{}

	err = xml.Unmarshal(body, &deviceList)
	if err != nil {
		return nil, err
	}

	servers := []*Server{}
	for _, device := range deviceList.Devices {
		// Device is a server and is owned by user
		if strings.Contains(device.Roles, "server") && device.Owned {
			// Loop over the server's connections and use the first one to work
			// If none of the connections work, the server is skipped
			for _, conn := range device.Connections {
				s, err := NewServer(config.PlexServerConfig{
					BaseURL:  conn.URI,
					Token:    device.AccessToken,
					Insecure: false,
				})
				if err != nil {
					fmt.Println(err)
					continue
				}
				servers = append(servers, s)
				break
			}
		}
	}

	return servers, nil
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
