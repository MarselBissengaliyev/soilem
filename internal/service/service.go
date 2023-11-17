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
}

type Session interface {
	CreateSession(session *model.Session) string
	RemoveSession(token string)
	GetSession(token string) (*model.Session, bool)
	GetUserName(token string) (model.UserName, bool)
	GetUserAgent(token string) (string, bool)
}

type Twilo interface {
	SendSMSConfirmation(*model.User) *model.Fail
}

type Service struct {
	User
	Session
	Twilo
}

func NewService(r *repo.Repository, cfg *configs.Config) *Service {
	return &Service{
		User:    NewUserService(r),
		Session: NewSessionService(),
		Twilo:   NewTwiloService(cfg),
	}
}
