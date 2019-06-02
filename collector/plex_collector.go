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
			"Number of active Plex sessions",
			[]string{"server"}, nil,
		),
	}
}

func (c *PlexCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.sessionsMetric
}

func (c *PlexCollector) Collect(ch chan<- prometheus.Metric) {
	serverMetrics := c.client.GetServerMetrics()

	for k, v := range serverMetrics {
		c.Logger.Trace(v)
		ch <- prometheus.MustNewConstMetric(c.sessionsMetric, prometheus.CounterValue, float64(v.ActiveSessions), k)
	}
}
