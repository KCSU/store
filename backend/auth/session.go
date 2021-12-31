package auth

import (
	"encoding/json"
	"time"
)

type Session struct {
	AuthURL      string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	IDToken      string
}

func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}
