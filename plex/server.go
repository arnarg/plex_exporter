package plex

import (
	"encoding/json"
	"fmt"

	"github.com/arnarg/plex_exporter/plex/api"
)

type Server struct {
	ID          string
	Name        string
	Product     string
	Version     string
	Platform    string
	AccessToken string
	URL         string
}

const TestURI = "%s/identity"
const StatusURI = "%s/status/sessions"

func NewServer(d *api.Device) (*Server, error) {
	server := &Server{
		ID:          d.ClientID,
		Name:        d.Name,
		Product:     d.Product,
		Version:     d.Version,
		Platform:    d.Platform,
		AccessToken: d.AccessToken,
	}

	// Pick a connection
	// First one to respond gets picked
	for _, conn := range d.Connections {
		url := fmt.Sprintf(TestURI, conn.URI)
		res, _, err := SendRequest("HEAD", url, map[string]string{})
		if err != nil {
			continue
		}

		if res.StatusCode == 200 {
			server.URL = conn.URI
			break
		}
	}

	if server.URL == "" {
		return nil, fmt.Errorf("Could not find a working connection for server \"%s\"", server.Name)
	}

	return server, nil
}

func (s *Server) GetSessionCount(h map[string]string) (int, error) {
	sessionListWrapper := api.SessionListWrapper{}

	url := fmt.Sprintf(StatusURI, s.URL)
	_, body, err := SendRequest("GET", url, AddTokenHeader(h, s.AccessToken))
	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(body, &sessionListWrapper)
	if err != nil {
		return -1, err
	}

	return sessionListWrapper.List.Size, nil
}
