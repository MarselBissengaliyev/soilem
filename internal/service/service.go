package service

import (
	"github.com/MarselBissengaliyev/soilem/configs"
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/repo"
)

type User interface {
	Registration(user *model.User) (*model.User, *model.Fail)
	Login(user *model.User) (*model.User, *model.Fail)
	GetUserByUserName(userName model.UserName) (*model.User, *model.Fail)
	GetUsers(searchTerm string, limit string) ([]*model.User, *model.Fail)
	ConfirmSMSCode(userName model.UserName, providedCode model.SMSCode) (bool, *model.Fail)
}

type Session interface {
	CreateSession(session *model.Session) string
	RemoveSession(token string)
	GetSession(token string) (*model.Session, bool)
	GetUserName(token string) (model.UserName, bool)
	GetUserAgent(token string) (string, bool)
}

type SMSCode interface {
	SendSMSConfirmation(*model.User) *model.Fail
	SetSMSCode(updateSMSCode model.SMSCode, userName model.UserName) (*model.SMSCode, *model.Fail)
}

type Service struct {
	User
	Session
	SMSCode
}

func NewService(r *repo.Repository, cfg *configs.Config) *Service {
	return &Service{
		User:    NewUserService(r),
		Session: NewSessionService(),
		SMSCode: NewSMSCodeService(cfg, r),
	}
}
