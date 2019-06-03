package api

type LibraryWrapper struct {
	Library Library `json:"MediaContainer"`
}

type Library struct {
	Size     int       `json:"size"`
	Sections []Section `json:"Directory"`
}

type Section struct {
	ID   string `json:"key"`
	Name string `json:"title"`
	Type string `json:"type"`
}
