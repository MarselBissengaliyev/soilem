package model

import (
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserPhone string
type UserName string
type UserID uint

type User struct {
	ID       UserID   `json:"-"`
	UserName UserName `json:"user_name" validate:"required"`

	// Phone
	PhoneNumber     UserPhone `json:"phone_number"`
	SMSCode         SMSCode   `json:"confirmation_code"`
	IsPhoneVerified bool      `json:"-"`
	// End Phone

	// Email
	Email           string    `json:"email" validate:"email"`
	EmailCode       EmailCode `json:"email_code" validate:"email_code"`
	IsEmailVerified bool      `json:"-"`
	// End Email

	Password string `json:"-" validate:"required"`

	// Profile
	Profile Profile `json:"profile"`
	// End Profile

	CreatedAt      time.Time `json:"-"`
	LastLoginAt    time.Time `json:"-"`
	IsRegistration bool      `json:"-"`
}

func (u *User) Validate() error {
	err := v.RegisterValidation("phone_number", func(fl validator.FieldLevel) bool {
		e164Regex := `^\+[1-9]\d{1,14}$`
		re := regexp.MustCompile(e164Regex)
		phone_number := strings.ReplaceAll(string(u.PhoneNumber), " ", "")

		return re.Find([]byte(phone_number)) != nil
	})
	if err != nil {
		return err
	}

	err = v.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		_, err := mail.ParseAddress(u.Email)
		return err == nil
	})
	if err != nil {
		return err
	}

	return v.Struct(u)
}

func (u *User) HashPassword() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	return string(bytes), err
}

func (u *User) CheckPasswordHash(provivedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(provivedPassword))
	return err == nil
}
