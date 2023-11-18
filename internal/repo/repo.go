package repo

import (
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/jackc/pgx/v5"
)

type User interface {
	CreateUser(user *model.User) (*model.User, error)
	GetUserByUserName(userName model.UserName) (*model.User, error)
	GetUsers(searchTerm string, limit int) ([]*model.User, error)
	SetPhoneVerifiedValue(status bool, userName model.UserName) (bool, error)
}

type SMSCode interface {
	SetSMSCode(updateSMSCode model.SMSCode, userName model.UserName) (*model.SMSCode, error)
}

type EmailCode interface {
	SetEmailCode(updateEmailCode model.EmailCode, userName model.UserName) (*model.EmailCode, error)
}

type Repository struct {
	User
	SMSCode
	EmailCode
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		User:      NewUserPostgres(db),
		SMSCode:   NewSMSCodePostgres(db),
		EmailCode: NewEmailCodePostgres(db),
	}
}
