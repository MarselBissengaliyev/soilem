package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RefreshToken struct {
	ID        uint      `json:"-"`
	Token     string    `json:"token"`
	UserName  string    `json:"user_name"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *RefreshToken) HashToken() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(t.Token), bcrypt.DefaultCost)
	return string(bytes), err
}

func (t *RefreshToken) CheckTokenHash(providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.Token), []byte(providedPassword))
	return err == nil
}

func (t *RefreshToken) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now().UTC())
}
