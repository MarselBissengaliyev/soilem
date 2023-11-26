package model

import (
	"math/rand"
	"time"
)

type SMSCode struct {
	ID        uint      `json:"-"`
	Code      int       `json:"code" validate:"required,len:6,numeric"`
	ExpiresAt time.Time `json:"-"`
	UserName  UserName  `json:"user_name"`
}

func (c *SMSCode) IsExpired() bool {
	return c.ExpiresAt.Before(time.Now().UTC())
}

func (c *SMSCode) Validate() error {
	return v.Struct(c)
}

func (c *SMSCode) GenerateConfirmationCode(userName UserName) {
	// Создаем новый источник случайных чисел с seed на основе текущего времени
	source := rand.NewSource(time.Now().UnixNano())

	// Создаем новый генератор случайных чисел на основе источника
	randomGenerator := rand.New(source)

	// Генерируем 6 случайных цифр
	randomNumber := randomGenerator.Intn(1000000)

	c.Code = randomNumber
	c.ExpiresAt = time.Now().Add(5 * time.Minute)
	c.UserName = userName
}

func (c *SMSCode) CheckSMSCode(providedCode int) bool {
	return c.Code == providedCode
}
