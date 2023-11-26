package repo

import (
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/jackc/pgx/v5"
)

type User interface {
	Create(user *model.User) (*model.User, pgx.Tx, error)
	GetByUserName(userName model.UserName) (*model.User, error)
	GetAll(searchTerm string, limit int) ([]*model.User, error)
	SetPhoneVerifiedValue(status bool, userName model.UserName) (bool, error)
	SetEmailVerifiedValue(status bool, userName model.UserName) (bool, error)
}

type SMSCode interface {
	SetCode(updateSMSCode model.SMSCode) (*model.SMSCode, error)
}

type EmailCode interface {
	SetCode(updateEmailCode model.EmailCode) (*model.EmailCode, error)
}

type Profile interface {
	Create(profile *model.Profile) (*model.Profile, error)
	GetByUserName(userName model.UserName) (*model.Profile, error)
}

type Post interface {
	Create(post *model.Post) (*model.Post, error)
	GetBySlug(slug string) (*model.Post, error)
}

type Repository struct {
	User
	SMSCode
	EmailCode
	Post
	Profile
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		User:      NewUserPostgres(db),
		SMSCode:   NewSMSCodePostgres(db),
		EmailCode: NewEmailCodePostgres(db),
		Profile:   NewProfilePostgres(db),
		Post:      NewPostPostgres(db),
	}
}
