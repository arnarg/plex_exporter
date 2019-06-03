package plex

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/arnarg/plex_exporter/plex/api"
	"github.com/imdario/mergo"
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
const LibraryURI = "%s/library/sections"
const SectionURI = "%s/library/sections/%d/all"

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

// GetLibrary returns server's library
func (s *Server) GetLibrary(h map[string]string) (*api.Library, error) {
	libraryWrapper := api.LibraryWrapper{}

	url := fmt.Sprintf(LibraryURI, s.URL)
	_, body, err := SendRequest("GET", url, AddTokenHeader(h, s.AccessToken))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &libraryWrapper)
	if err != nil {
		return nil, err
	}

	return &libraryWrapper.Library, err
}

// GetSectionSize returns the number of items in a library section with
// given ID.
func (s *Server) GetSectionSize(i int, h map[string]string) (int, error) {
	// If certain headers are added to this request it returns the size of
	// the library section as a header. Therefor we can just make a HEAD request.
	eh := map[string]string{
		"X-Plex-Container-Start": "0",
		"X-Plex-Container-Size":  "0",
		"X-Plex-Sync-Version":    "2",
	}
	mergo.Merge(&eh, h)

	url := fmt.Sprintf(SectionURI, s.URL, i)
	resp, _, err := SendRequest("HEAD", url, AddTokenHeader(eh, s.AccessToken))
	if err != nil {
		return -1, err
	}

	size, ok := resp.Header["X-Plex-Container-Total-Size"]
	if !ok {
		return -1, fmt.Errorf("Did not receive X-Plex-Container-Total-Size header")
	}

	ret, err := strconv.Atoi(size[0])
	if !ok {
		return -1, fmt.Errorf("Could not parse size as int")
	}

	return ret, nil
}
