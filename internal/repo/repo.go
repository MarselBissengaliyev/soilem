package repo

import (
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/jackc/pgx/v5"
)

type User interface {
	Create(user *model.User) (*model.User, error)
	GetByUserName(userName model.UserName) (*model.User, error)
	GetAll(searchTerm string, limit int) ([]*model.User, error)
	SetPhoneVerifiedValue(status bool, userName model.UserName) (bool, error)
	SetEmailVerifiedValue(status bool, userName model.UserName) (bool, error)
}

type SMSCode interface {
	SetCode(updateSMSCode model.SMSCode, userName model.UserName) (*model.SMSCode, error)
}

type EmailCode interface {
	SetCode(updateEmailCode model.EmailCode, userName model.UserName) (*model.EmailCode, error)
}

type Post interface {
	
}

type Repository struct {
	User
	SMSCode
	EmailCode
	Post
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		User:      NewUserPostgres(db),
		SMSCode:   NewSMSCodePostgres(db),
		EmailCode: NewEmailCodePostgres(db),
		Post:      NewPostPostgres(db),
	}
}
