package service

import (
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/pkg/utils"
)

type SessionService struct {
	sessions map[string]model.Session
}

func NewSessionService() *SessionService {
	return &SessionService{sessions: make(map[string]model.Session)}
}

func (s *SessionService) CreateSession(session *model.Session) string {
	token := utils.GenerateUniqueToken()

	s.sessions[token] = model.Session{
		UserName:  session.UserName,
		Expiry:    session.Expiry,
		UserAgent: session.UserAgent,
	}

	return token
}

func (s *SessionService) GetUserName(token string) (model.UserName, bool) {
	sess, ok := s.GetSession(token)
	if !ok {
		return "", false
	}

	return sess.UserName, true
}

func (s *SessionService) GetUserAgent(token string) (string, bool) {
	sess, ok := s.GetSession(token)
	if !ok {
		return "", false
	}

	return sess.UserAgent, true
}

func (s *SessionService) RemoveSession(token string) {
	delete(s.sessions, token)
}

func (s *SessionService) GetSession(token string) (*model.Session, bool) {
	sess, exists := s.sessions[token]
	if !exists {
		return nil, false
	}

	if sess.IsExpired() {
		s.RemoveSession(token)
		return nil, false
	}

	return &sess, true
}
