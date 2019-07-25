package api

type SessionListWrapper struct {
	List SessionList `json:"MediaContainer"`
}

type SessionList struct {
	Size     int       `json:"size"`
	Sessions []Session `json:"Metadata"`
}

type Session struct {
	Player Player `json:"Player"`
}

type Player struct {
	Address   string `json:"address"`
	MachineID string `json:"machineIdentifier"`
}
