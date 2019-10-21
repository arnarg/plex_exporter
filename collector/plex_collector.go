package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/tagnard/plex_exporter/plex"
)

type PlexCollector struct {
	Logger         *log.Entry
	client         *plex.PlexClient
	serverInfo     *prometheus.Desc
	sessionsMetric *prometheus.Desc
	libraryMetric  *prometheus.Desc
}

func NewPlexCollector(c *plex.PlexClient, l *log.Entry) *PlexCollector {
	return &PlexCollector{
		Logger: l,
		client: c,
		serverInfo: prometheus.NewDesc("plex_server_info",
			"Information about Plex server",
			[]string{"server_name", "server_id", "version", "platform"}, nil,
		),
		sessionsMetric: prometheus.NewDesc("plex_sessions_total",
			"Number of active Plex sessions",
			[]string{"server_name", "server_id", "username", "library", "state", "title"}, nil,
		),
		libraryMetric: prometheus.NewDesc("plex_library_section_size_count",
			"Number of items in a library section",
			[]string{"server_name", "server_id", "name", "type"}, nil,
		),
	}
}

func (c *PlexCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.serverInfo
	ch <- c.sessionsMetric
	ch <- c.libraryMetric
}

func (c *PlexCollector) Collect(ch chan<- prometheus.Metric) {
	serverMetrics := c.client.GetServerMetrics()

	for _, v := range serverMetrics {
		c.Logger.Trace(v)
		ch <- prometheus.MustNewConstMetric(c.serverInfo, prometheus.CounterValue, 1, v.Name, v.ID, v.Version, v.Platform)

		for _, l := range v.Sessions {
			ch <- prometheus.MustNewConstMetric(c.sessionsMetric, prometheus.GaugeValue, float64(1), v.Name, v.ID, l.Username, l.Library, l.State, l.Title)
		}

		for _, l := range v.Libraries {
			ch <- prometheus.MustNewConstMetric(c.libraryMetric, prometheus.GaugeValue, float64(l.Size), v.Name, v.ID, l.Name, l.Type)
		}
	}
}
