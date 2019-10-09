package api

type SessionList struct {
	Sessions `json:"MediaContainer"`
}

type Sessions struct {
	Size int `json:"size"`
}
