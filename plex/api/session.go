package api

type SessionList struct {
	Sessions `json:"MediaContainer"`
}

type Sessions struct {
	Metadata []Metadata `json:"Metadata"`
}

type User struct {
	Username string `json:"title"`
}

type Player struct {
	State string `json:"state"`
}

type Metadata struct {
	User    `json:"User"`
	Player  `json:"Player"`
	Library string `json:"librarySectionTitle"`
	Series  string `json:"grandparentTitle"`
	Season  string `json:"parentTitle"`
	Title   string `json:"title"`
	Type    string `json:"type"`
}
