package plex

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"runtime"
	"strings"

	"github.com/arnarg/plex_exporter/config"
	"github.com/arnarg/plex_exporter/plex/api"
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
)

type PlexClient struct {
	Logger         *log.Entry
	Token          string
	Servers        []api.Device
	DefaultHeaders map[string]string
}

func NewPlexClient(c *config.PlexConfig, v string, l *log.Entry) *PlexClient {
	return &PlexClient{
		Logger:  l,
		Token:   c.Token,
		Servers: []api.Device{},
		DefaultHeaders: map[string]string{
			"User-Agent":               fmt.Sprintf("plex_exporter/%s", v),
			"Accept":                   "application/json",
			"X-Plex-Platform":          runtime.GOOS,
			"X-Plex-Version":           v,
			"X-Plex-Client-Identifier": c.UUID,
			"X-Plex-Device-Name":       "Plex Exporter",
		},
	}
}

// Init fetches all Plex Media Servers and stores a reference to them
func (c *PlexClient) Init() error {
	if c.Token == "" {
		return fmt.Errorf("Authentication token missing")
	}

	// This endpoint only supports XML.
	// I want to specify the "Accept: application/xml" header
	// to make sure that if the endpoint does support JSON in
	// the future it won't break the application.
	h := map[string]string{
		"Accept": "application/xml",
	}
	mergo.Merge(&h, c.DefaultHeaders)

	_, body, err := SendRequest("GET", "https://plex.tv/api/resources?includeHttps=1", AddTokenHeader(h, c.Token))
	if err != nil {
		return err
	}

	deviceList := api.DeviceList{}

	err = xml.Unmarshal(body, &deviceList)
	if err != nil {
		return err
	}

	for _, device := range deviceList.Devices {
		if strings.Contains(device.Roles, "server") {
			c.Logger.Debugf("Found server %s", device.Name)
			c.Servers = append(c.Servers, device)
		}
	}

	return nil
}

// GetSessions fetches sessions from all the servers and stores them in a map
// using the server's name as the key.
func (c *PlexClient) GetSessions() (*map[string]api.SessionList, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("Authentication token missing")
	}

	serverMap := map[string]api.SessionList{}

	// Fetch session list from servers
	// I don't love this approach but it works for now
	for _, server := range c.Servers {
		logger := c.Logger.WithFields(log.Fields{"server": server.Name})
		sessionListWrapper := api.SessionListWrapper{}
		successful := false

		for _, conn := range server.Connections {
			logger.Debugf("Getting session list from URL %s", conn.URI)

			url := fmt.Sprintf("%s/status/sessions", conn.URI)
			_, body, err := SendRequest("GET", url, AddTokenHeader(c.DefaultHeaders, c.Token))
			if err != nil {
				logger.Debugf("Couldn't fetch session list from URL %s: %s", conn.URI, err)
				continue
			}

			err = json.Unmarshal(body, &sessionListWrapper)
			if err != nil {
				logger.Debugf("Couldn't parse session list response from URL %s: %s", conn.URI, err)
				logger.Trace(string(body))
			}

			successful = true
			break
		}

		if successful {
			c.Logger.Debugf("Successfully got session list from server \"%s\"", server.Name)
			serverMap[server.Name] = sessionListWrapper.List
		} else {
			c.Logger.Errorf("Unable to get session list from server \"%s\"", server.Name)
		}
	}
	return &serverMap, nil
}

// GetPinRequest creates a PinRequest using the Plex API and returns it.
func (c *PlexClient) GetPinRequest() (*api.PinRequest, error) {
	_, body, err := SendRequest("POST", "https://plex.tv/pins", c.DefaultHeaders)
	if err != nil {
		return nil, err
	}

	container := &api.PinRequestContainer{}

	err = json.Unmarshal(body, container)
	if err != nil {
		return nil, err
	}

	return &container.PinRequest, nil
}

// GetTokenFromPinRequest takes in a PinRequest and checks if it has been authenticated.
// If it has been authenticated it returns the token.
// If it has not been authenticated it returns an empty string.
func (c *PlexClient) GetTokenFromPinRequest(p *api.PinRequest) (string, error) {
	_, body, err := SendRequest("GET", fmt.Sprintf("https://plex.tv/pins/%d", p.Id), c.DefaultHeaders)
	if err != nil {
		return "", err
	}

	container := api.PinRequestContainer{}

	err = json.Unmarshal(body, &container)
	if err != nil {
		return "", err
	}

	if container.PinRequest.AuthToken != "" {
		c.Token = container.PinRequest.AuthToken
	}

	return c.Token, nil
}
