package plex

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/arnarg/plex_exporter/config"
	"github.com/arnarg/plex_exporter/plex/api"
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
)

type PlexClient struct {
	Logger         *log.Entry
	Token          string
	Servers        []Server
	DefaultHeaders map[string]string
}

func NewPlexClient(c *config.PlexConfig, v string, l *log.Entry) *PlexClient {
	return &PlexClient{
		Logger:  l,
		Token:   c.Token,
		Servers: []Server{},
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
			c.Logger.Debugf("Found server \"%s\"", device.Name)
			server, err := NewServer(&device)
			if err != nil {
				c.Logger.Errorf("Could not use server: %s", err)
				continue
			}
			c.Servers = append(c.Servers, *server)
		}
	}

	if len(c.Servers) < 1 {
		return fmt.Errorf("No suitable servers found.")
	}

	return nil
}

// GetServerMetrics fetches all metrics for each server and returns them in a map
// with the servers' names as keys.
func (c *PlexClient) GetServerMetrics() map[string]ServerMetric {
	serverMap := map[string]ServerMetric{}

	for _, server := range c.Servers {
		logger := c.Logger.WithFields(log.Fields{"server": server.Name})

		serverMetric := ServerMetric{
			ID:       server.ID,
			Name:     server.Name,
			Product:  server.Product,
			Version:  server.Version,
			Platform: server.Platform,
		}

		// Get active sessions
		activeSessions, err := server.GetSessionCount(c.DefaultHeaders)
		if err != nil {
			logger.Errorf("Could not get metrics for server \"%s\"", server.Name)
			logger.Debugf("Could not get session count: %s", err)
			continue
		}
		serverMetric.ActiveSessions = activeSessions

		// Get library metrics
		library, err := server.GetLibrary(c.DefaultHeaders)
		if err != nil {
			logger.Errorf("Could not get metrics for server \"%s\"", server.Name)
			logger.Debugf("Could not get library: %s", err)
			continue
		}

		for _, section := range library.Sections {
			id, err := strconv.Atoi(section.ID)
			if err != nil {
				logger.Debugf("Could not convert sections ID to int. (%s)", section.ID)
			}
			size, err := server.GetSectionSize(id, c.DefaultHeaders)
			if err != nil {
				logger.Debugf("Could not get section size for \"%s\": %s", section.Name, err)
				continue
			}
			libraryMetric := LibraryMetric{
				Name: section.Name,
				Type: section.Type,
				Size: size,
			}

			serverMetric.Libraries = append(serverMetric.Libraries, libraryMetric)
		}

		serverMap[server.Name] = serverMetric
	}

	return serverMap
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
