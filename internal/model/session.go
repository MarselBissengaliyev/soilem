package model

import (
	"time"
)

type Session struct {
	UserName  UserName
	UserAgent string
	Expiry    time.Time
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now().UTC())
}
