package plex

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

type PlexClient struct {
	Logger *log.Entry
	server *Server
}

func NewPlexClient(s *Server, l *log.Entry) (*PlexClient, error) {
	return &PlexClient{
		Logger: l,
		server: s,
	}, nil
}

// GetServerMetrics fetches all metrics for each server and returns them in a map
// with the servers' names as keys.
func (c *PlexClient) GetServerMetrics() ServerMetric {
	logger := c.Logger.WithFields(log.Fields{"server": c.server.Name})

	serverMetric := ServerMetric{
		Version:  c.server.Version,
		Platform: c.server.Platform,
	}

	// Get active sessions
	activeSessions, err := c.server.GetSessionCount()
	if err != nil {
		logger.Errorf("Could not get metrics")
		logger.Debugf("Could not get session count: %s", err)
		// TODO fix
		return serverMetric
	}
	serverMetric.ActiveSessions = activeSessions

	// Get library metrics
	library, err := c.server.GetLibrary()
	if err != nil {
		logger.Errorf("Could not get metrics")
		logger.Debugf("Could not get library: %s", err)
		return serverMetric
	}

	for _, section := range library.Sections {
		id, err := strconv.Atoi(section.ID)
		if err != nil {
			logger.Debugf("Could not convert sections ID to int. (%s)", section.ID)
		}
		size, err := c.server.GetSectionSize(id)
		if err != nil {
			logger.Debugf("Could not get section size for \"%s\": %s", section.Name, err)
			return serverMetric
		}
		libraryMetric := LibraryMetric{
			Name: section.Name,
			Type: section.Type,
			Size: size,
		}

		serverMetric.Libraries = append(serverMetric.Libraries, libraryMetric)
	}

	return serverMetric
}
