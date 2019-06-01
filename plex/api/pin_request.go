package api

import (
	"time"
)

// PinRequest represents a pin request
type PinRequest struct {
	Id        int       `json:"id"`
	Code      string    `json:"code"`
	Expiry    time.Time `json:"expires_at"`
	Trusted   bool      `json:"trusted"`
	AuthToken string    `json:"auth_token"`
}

// PinContainer is a wrapper around a pin request.
// The API returns {"pin": {"id": 1,...}} so this is required for parsing.
type PinRequestContainer struct {
	PinRequest PinRequest `json:"pin"`
}
