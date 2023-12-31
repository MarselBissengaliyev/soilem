package service

import (
	"github.com/MarselBissengaliyev/soilem/configs"
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/repo"
	"github.com/jackc/pgx/v5"
)

type User interface {
	Registration(user *model.User) (*model.User, pgx.Tx, *model.Fail)
	Login(user *model.User) (*model.User, *model.Fail)
	GetByUserName(userName model.UserName) (*model.User, *model.Fail)
	GetUsers(searchTerm string, limit string) ([]*model.User, *model.Fail)
	ConfirmSMSCode(userName model.UserName, providedCode model.SMSCode) (bool, *model.Fail)
	ConfirmEmailCode(userName model.UserName, providedCode model.EmailCode) (bool, *model.Fail)
}

type AccessToken interface {
	Create(session *model.AccessToken) (string, *model.Fail)
	RemoveByAccessToken(token string) *model.Fail
	GetByAccessToken(token string) (*model.Session, *model.Fail)
}

type SMSCode interface {
	SendSMSConfirmation(to model.UserPhone, code int) *model.Fail
	SetSMSCode(updateSMSCode model.SMSCode, userName model.UserName) (*model.SMSCode, *model.Fail)
}

type EmailCode interface {
	SetEmailCode(updateEmailCode model.EmailCode, userName model.UserName) (*model.EmailCode, *model.Fail)
	SendEmailCode(templatePath string, to string, code int) *model.Fail
}

type Profile interface {
	Create(profile *model.Profile) (*model.Profile, *model.Fail)
}

type Post interface {
	Create(post *model.Post, userName model.UserName) (*model.Post, *model.Fail)
}

type Service struct {
	User
	AccessToken
	SMSCode
	EmailCode
	Profile
	Post
}

func NewService(r *repo.Repository, cfg *configs.Config) *Service {
	return &Service{
		User:        NewUserService(r.User),
		AccessToken: NewAccessTokenService(r.AccessToken),
		SMSCode:     NewSMSCodeService(cfg, r.SMSCode),
		EmailCode:   NewEmailCodeService(cfg, r.EmailCode),
		Profile:     NewProfileService(r.Profile),
		Post:        NewPostService(r.Post),
	}
}
