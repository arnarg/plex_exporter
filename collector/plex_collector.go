package collector

import (
	"github.com/arnarg/plex_exporter/plex"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type PlexCollector struct {
	Logger         *log.Entry
	client         *plex.PlexClient
	sessionsMetric *prometheus.Desc
}

func NewPlexCollector(c *plex.PlexClient, l *log.Entry) *PlexCollector {
	return &PlexCollector{
		Logger: l,
		client: c,
		sessionsMetric: prometheus.NewDesc("plex_active_sessions_count",
			"Shows the number of active sessions",
			[]string{"server"}, nil,
		),
	}
}

func (c *PlexCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.sessionsMetric
}

func (c *PlexCollector) Collect(ch chan<- prometheus.Metric) {
	sessionLists, err := c.client.GetSessions()
	if err != nil {
		c.Logger.Errorf("Could not get session lists: %s", err)
	}

	for key, value := range *sessionLists {
		c.Logger.Trace(value)
		ch <- prometheus.MustNewConstMetric(c.sessionsMetric, prometheus.CounterValue, float64(value.Size), key)
	}
}
