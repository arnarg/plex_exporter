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
			[]string{"version", "platform"},
		),
		sessionsMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "plex",
				Subsystem: "sessions",
				Name:      "active_count",
				Help:      "Number of active Plex sessions",
			},
			[]string{},
		),
		libraryMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "plex",
				Subsystem: "library",
				Name:      "section_size_count",
				Help:      "Number of items in a library section",
			},
			[]string{"name", "type"},
		),
	}
}

func (c *PlexCollector) Describe(ch chan<- *prometheus.Desc) {
	c.serverInfo.Describe(ch)
	c.sessionsMetric.Describe(ch)
	c.libraryMetric.Describe(ch)
}

func (c *PlexCollector) Collect(ch chan<- prometheus.Metric) {
	v, err := c.client.GetServerMetrics()
	if err != nil {
		c.Logger.Errorf("Could not retrieve server metrics: %s", err)
		return
	}

	c.Logger.Trace(v)
	c.serverInfo.WithLabelValues(v.Version, v.Platform).Set(1)
	c.sessionsMetric.WithLabelValues().Set(float64(v.ActiveSessions))

	for _, l := range v.Libraries {
		c.libraryMetric.WithLabelValues(l.Name, l.Type).Set(float64(l.Size))
	}

	c.serverInfo.Collect(ch)
	c.sessionsMetric.Collect(ch)
	c.libraryMetric.Collect(ch)
}
