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

type ShowSectionResponse struct {
	ShowSectionDetail `json:"MediaContainer"`
}

type ShowSectionDetail struct {
	ShowCount int          `json:"size"`
	Shows     []ShowDetail `json:"MetaData"`
}

type ShowDetail struct {
	EpisodeCount        int `json:"leafCount"`
	WatchedEpisodeCount int `json:"viewedLeafCount"`
	SeasonCount         int `json:"childCount"`
}
