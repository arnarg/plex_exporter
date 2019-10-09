package api

type LibraryResponse struct {
	Library `json:"MediaContainer"`
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

type SectionResponse struct {
	SectionDetail `json:"MediaContainer"`
}

type SectionDetail struct {
	TotalSize int `json:"totalSize"`
}
