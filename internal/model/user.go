package model

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var v = validator.New()

type UserPhone string
type UserName string
type UserID uint

type confirmationCode struct {
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

type User struct {
	ID UserID `json:"-"`

	// Phone
	PhoneNumber      UserPhone        `json:"phone_number"`
	ConfirmationCode confirmationCode `json:"confirmation_code"`
	IsPhoneVerified  bool             `json:"-"`
	// End Phone

	Password       string    `json:"-" validate:"required"`
	FullName       string    `json:"full_name" validate:"required"`
	UserName       UserName  `json:"user_name" validate:"required"`
	Description    string    `json:"description"`
	Avatar         string    `json:"avatar"`
	DateOfBirth    time.Time `json:"birth_of_date" validate:"required"`
	Sex            string    `json:"sex" validate:"required, oneof=male female"`
	CreatedAt      time.Time `json:"-"`
	LastLoginAt    time.Time `json:"-"`
	IsRegistration bool      `json:"-"`
}

func (u *User) Validate() error {
	if u.IsRegistration {
		if err := v.Var(u.PhoneNumber, "required"); err != nil {
			return err
		}

		if err := v.Var(u.FullName, "required"); err != nil {
			return err
		}
	}

	if err := v.RegisterValidation("phone_number", func(fl validator.FieldLevel) bool {
		e164Regex := `^\+[1-9]\d{1,14}$`
		re := regexp.MustCompile(e164Regex)
		phone_number := strings.ReplaceAll(string(u.PhoneNumber), " ", "")

		return re.Find([]byte(phone_number)) != nil
	}); err != nil {
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

func (u *User) GenerateConfirmationCode() {
	// Создаем новый источник случайных чисел с seed на основе текущего времени
	source := rand.NewSource(time.Now().UnixNano())

	// Создаем новый генератор случайных чисел на основе источника
	randomGenerator := rand.New(source)

	// Генерируем 6 случайных цифр
	randomNumber := randomGenerator.Intn(1000000)

	// Форматируем число как строку с шестью цифрами
	randomString := fmt.Sprintf("%06d", randomNumber)

	u.ConfirmationCode = confirmationCode{
		Code:      randomString,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}
