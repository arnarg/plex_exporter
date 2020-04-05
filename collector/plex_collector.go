package collector

import (
	"github.com/arnarg/plex_exporter/plex"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type PlexCollector struct {
	Logger *log.Entry
	client *plex.PlexClient

	serverInfo     *prometheus.GaugeVec
	sessionsMetric *prometheus.GaugeVec
	libraryMetric  *prometheus.GaugeVec
}

func NewPlexCollector(c *plex.PlexClient, l *log.Entry) *PlexCollector {
	return &PlexCollector{
		Logger: l,
		client: c,

		serverInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "plex",
				Subsystem: "server",
				Name:      "info",
				Help:      "Information about Plex server",
			},
			[]string{"server_name", "server_id", "version", "platform"},
		),
		sessionsMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "plex",
				Subsystem: "sessions",
				Name:      "active_count",
				Help:      "Number of active Plex sessions",
			},
			[]string{"server_name", "server_id"},
		),
		libraryMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "plex",
				Subsystem: "library",
				Name:      "section_size_count",
				Help:      "Number of items in a library section",
			},
			[]string{"server_name", "server_id", "name", "type"},
		),
	}
}

func (c *PlexCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *PlexCollector) Collect(ch chan<- prometheus.Metric) {
	serverMetrics := c.client.GetServerMetrics()

	for _, v := range serverMetrics {
		c.Logger.Trace(v)
		c.serverInfo.WithLabelValues(v.Name, v.ID, v.Version, v.Platform).Set(1)
		c.sessionsMetric.WithLabelValues(v.Name, v.ID).Set(float64(v.ActiveSessions))

		for _, l := range v.Libraries {
			c.libraryMetric.WithLabelValues(v.Name, v.ID, l.Name, l.Type).Set(float64(l.Size))
		}
	}

	c.serverInfo.Collect(ch)
	c.sessionsMetric.Collect(ch)
	c.libraryMetric.Collect(ch)
}
