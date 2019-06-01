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
	User   User   `json:"User"`
}

type Player struct {
	Address   string `json:"address"`
	MachineID string `json:"machineIdentifier"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"title"`
}
