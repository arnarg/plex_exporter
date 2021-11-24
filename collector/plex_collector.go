package collector

import (
	"github.com/arnarg/plex_exporter/plex"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type PlexCollector struct {
	Logger                   *log.Entry
	client                   *plex.PlexClient
	serverInfo               *prometheus.Desc
	sessionsMetric           *prometheus.Desc
	libraryMetric            *prometheus.Desc
	showLibraryMetric        *prometheus.Desc
	showLibrarySeasonMetric  *prometheus.Desc
	showLibraryEpisodeMetric *prometheus.Desc
	showLibraryWatchedMetric *prometheus.Desc
}

func NewPlexCollector(c *plex.PlexClient, l *log.Entry) *PlexCollector {
	return &PlexCollector{
		Logger: l,
		client: c,
		serverInfo: prometheus.NewDesc("plex_server_info",
			"Information about Plex server",
			[]string{"server_name", "server_id", "version", "platform"}, nil,
		),
		sessionsMetric: prometheus.NewDesc("plex_sessions_active_count",
			"Number of active Plex sessions",
			[]string{"server_name", "server_id"}, nil,
		),
		libraryMetric: prometheus.NewDesc("plex_library_section_size_count",
			"Number of items in a library section",
			[]string{"server_name", "server_id", "name", "type"}, nil,
		),
		showLibraryMetric: prometheus.NewDesc("plex_library_section_show_count",
			"Number of shows in a library section of type show",
			[]string{"server_name", "server_id", "name"}, nil,
		),
		showLibrarySeasonMetric: prometheus.NewDesc("plex_library_section_show_season_count",
			"Number of seasons in a library section of type show",
			[]string{"server_name", "server_id", "name"}, nil,
		),
		showLibraryEpisodeMetric: prometheus.NewDesc("plex_library_section_show_episode_count",
			"Number of episodes in a library section of type show",
			[]string{"server_name", "server_id", "name"}, nil,
		),
		showLibraryWatchedMetric: prometheus.NewDesc("plex_library_section_show_watched_count",
			"Number of watched episodes in a library section of type show",
			[]string{"server_name", "server_id", "name"}, nil,
		),
	}
}

func (c *PlexCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.serverInfo
	ch <- c.sessionsMetric
	ch <- c.libraryMetric
	ch <- c.showLibraryMetric
	ch <- c.showLibrarySeasonMetric
	ch <- c.showLibraryEpisodeMetric
	ch <- c.showLibraryWatchedMetric
}

func (c *PlexCollector) Collect(ch chan<- prometheus.Metric) {
	serverMetrics := c.client.GetServerMetrics()

	for _, v := range serverMetrics {
		c.Logger.Trace(v)
		ch <- prometheus.MustNewConstMetric(c.serverInfo, prometheus.CounterValue, 1, v.Name, v.ID, v.Version, v.Platform)
		ch <- prometheus.MustNewConstMetric(c.sessionsMetric, prometheus.GaugeValue, float64(v.ActiveSessions), v.Name, v.ID)

		for _, l := range v.Libraries {
			ch <- prometheus.MustNewConstMetric(c.libraryMetric, prometheus.GaugeValue, float64(l.Size), v.Name, v.ID, l.Name, l.Type)
		}
		for _, s := range v.ShowLibraries {
			ch <- prometheus.MustNewConstMetric(c.showLibraryMetric, prometheus.GaugeValue, float64(s.ShowSize), v.Name, v.ID, s.Name)
			ch <- prometheus.MustNewConstMetric(c.showLibrarySeasonMetric, prometheus.GaugeValue, float64(s.SeasonSize), v.Name, v.ID, s.Name)
			ch <- prometheus.MustNewConstMetric(c.showLibraryEpisodeMetric, prometheus.GaugeValue, float64(s.EpisodeSize), v.Name, v.ID, s.Name)
			ch <- prometheus.MustNewConstMetric(c.showLibraryWatchedMetric, prometheus.GaugeValue, float64(s.WatchedSize), v.Name, v.ID, s.Name)
		}
	}
}
