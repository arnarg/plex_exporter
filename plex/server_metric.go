package plex

type ServerMetric struct {
	Version        string
	Platform       string
	ActiveSessions int
	Libraries      []LibraryMetric
}

type LibraryMetric struct {
	Name string
	Type string
	Size int
}
