package plex

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/arnarg/plex_exporter/config"
	"github.com/arnarg/plex_exporter/plex/api"
	log "github.com/sirupsen/logrus"
)

type PlexClient struct {
	Logger  *log.Entry
	Token   string
	Servers []api.Device
	Headers map[string]string
}

func NewPlexClient(c *config.PlexConfig, v string, l *log.Entry) *PlexClient {
	return &PlexClient{
		Logger:  l,
		Token:   c.Token,
		Servers: []api.Device{},
		Headers: map[string]string{
			"User-Agent":               fmt.Sprintf("plex_exporter/%s", v),
			"Accept":                   "application/json",
			"X-Plex-Platform":          runtime.GOOS,
			"X-Plex-Version":           v,
			"X-Plex-Client-Identifier": c.UUID,
			"X-Plex-Device-Name":       "Plex Exporter",
		},
	}
}

func (c *PlexClient) Init() error {
	if c.Token == "" {
		return fmt.Errorf("Authentication token missing")
	}

	// This endpoint only supports XML.
	// I want to specify the "Accept: application/xml" header
	// to make sure that if the endpoint does support JSON in
	// the future it won't break the application.
	extraHeaders := &map[string]string{
		"Accept": "application/xml",
	}

	body, err := c.SendRequest("GET", "https://plex.tv/api/resources?includeHttps=1", extraHeaders)
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

func (c *PlexClient) GetSessions() (*map[string]api.SessionList, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("Authentication token missing")
	}

	serverMap := map[string]api.SessionList{}

	// Fetch session list from servers
	for _, server := range c.Servers {
		logger := c.Logger.WithFields(log.Fields{"server": server.Name})
		sessionListWrapper := api.SessionListWrapper{}
		successful := false

		for _, conn := range server.Connections {
			logger.Debugf("Getting session list from URL %s", conn.URI)

			url := fmt.Sprintf("%s/status/sessions", conn.URI)
			body, err := c.SendRequest("GET", url, nil)
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
	body, err := c.SendRequest("POST", "https://plex.tv/pins", nil)
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
	body, err := c.SendRequest("GET", fmt.Sprintf("https://plex.tv/pins/%d", p.Id), nil)
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

// SendRequest sends a HTTP request using a method parameter to a url parameter.
func (c *PlexClient) SendRequest(method, url string, h *map[string]string) ([]byte, error) {
	req, err := c.CreateRequest(method, url, h)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: time.Second * 10}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// CreateRequest creates a HTTP request including all headers from the PlexClient.
// If the PlexClient has a token that is included in a header as well.
func (c *PlexClient) CreateRequest(method, url string, h *map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range c.Headers {
		req.Header.Set(key, val)
	}

	if c.Token != "" {
		req.Header.Set("X-Plex-Token", c.Token)
	}

	// Overwrite with headers passed in parameter
	if h != nil {
		for key, val := range *h {
			req.Header.Set(key, val)
		}
	}

	return req, nil
}
