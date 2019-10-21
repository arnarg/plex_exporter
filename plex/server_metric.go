package plex

type ServerMetric struct {
	ID             string
	Name           string
	Version        string
	Platform       string
	ActiveSessions int
	Libraries      []LibraryMetric
	Sessions       []SessionMetric
}

type LibraryMetric struct {
	Name string
	Type string
	Size int
}

type SessionMetric struct {
	Username string
	Library  string
	State    string
	Title    string
}
