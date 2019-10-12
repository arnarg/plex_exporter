package plex

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/arnarg/plex_exporter/config"
	v "github.com/arnarg/plex_exporter/version"
	log "github.com/sirupsen/logrus"
)

var headers = map[string]string{
	"User-Agent":               fmt.Sprintf("plex_exporter/%s", v.Version),
	"Accept":                   "application/json",
	"X-Plex-Platform":          runtime.GOOS,
	"X-Plex-Version":           v.Version,
	"X-Plex-Client-Identifier": fmt.Sprintf("plex-exporter-v%s", v.Version),
	"X-Plex-Device-Name":       "Plex Exporter",
	"X-Plex-Product":           "Plex Exporter",
	"X-Plex-Device":            runtime.GOOS,
}

type PlexClient struct {
	Logger  *log.Entry
	Servers []*Server
	headers map[string]string
}

func NewPlexClient(c *config.PlexConfig, l *log.Entry) (*PlexClient, error) {
	var serverList []*Server

	h := headers
	h["X-Plex-Token"] = c.Token

	for _, serverConf := range c.Servers {
		plexServer, err := NewServer(serverConf)
		if err != nil {
			l.Errorf("Could not add server %s: %s", serverConf.BaseURL, err)
		} else {
			serverList = append(serverList, plexServer)
		}
	}

	if c.AutoDiscover {
		discoveryList, err := discoverServers(h)
		if err == nil {
			serverList = append(serverList, discoveryList...)
		}
	}

	l.Infof("Found %d working servers", len(serverList))

	return &PlexClient{
		Logger:  l,
		Servers: serverList,
		headers: h,
	}, nil
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
			Version:  server.Version,
			Platform: server.Platform,
		}

		// Get active sessions
		activeSessions, err := server.GetSessionCount()
		if err != nil {
			logger.Errorf("Could not get metrics for server \"%s\"", server.Name)
			logger.Debugf("Could not get session count: %s", err)
			continue
		}
		serverMetric.ActiveSessions = activeSessions

		// Get library metrics
		library, err := server.GetLibrary()
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
			size, err := server.GetSectionSize(id)
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
