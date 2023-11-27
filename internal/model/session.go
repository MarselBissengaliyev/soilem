package model

import (
	"time"
)

type Session struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	UserName     UserName  `json:"user_name"`
	UserAgent    string    `json:"user_agent"`
	ExpiresAt    time.Time `json:"expirest_at"`
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now().UTC())
}
