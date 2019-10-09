package api

type ServerInfoResponse struct {
	ServerInfo `json:"MediaContainer"`
}

type ServerInfo struct {
	ID       string `json:"machineIdentifier"`
	Name     string `json:"friendlyName"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
}
