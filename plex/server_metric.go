package plex

type ServerMetric struct {
	ID             string
	Name           string
	Version        string
	Platform       string
	ActiveSessions int
	Libraries      []LibraryMetric
	ShowLibraries  []ShowLibraryMetric
}

type LibraryMetric struct {
	Name string
	Type string
	Size int
}

type ShowLibraryMetric struct {
	Name        string
	Type        string
	ShowSize    int
	SeasonSize  int
	EpisodeSize int
	WatchedSize int
}
