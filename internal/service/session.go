package service

import (
	"net/http"

	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/repo"
	"github.com/MarselBissengaliyev/soilem/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type AccessTokenService struct {
	repo repo.AccessToken
}

func NewAccessTokenService(repo repo.AccessToken) *AccessTokenService {
	return &AccessTokenService{repo}
}

func (s *AccessTokenService) Create(accessToken *model.AccessToken) (string, *model.Fail) {
	token := utils.GenerateUniqueToken()

	accessToken.Token = token
	accessToken.HashToken()

	err := s.repo.Create(accessToken)
	if err != nil {
		return "", &model.Fail{
			Message:    "failed to create access token: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return token, nil
}

func (s *AccessTokenService) RemoveByAccessToken(accessToken string) *model.Fail {
	if accessToken == "" {
		return &model.Fail{
			Message:    "access token cannot be empty",
			StatusCode: http.StatusBadRequest,
		}
	}

	if err := s.repo.RemoveByToken(accessToken); err != nil {
		return &model.Fail{
			Message:    "failed to remove access token: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (s *AccessTokenService) GetByAccessToken(accessToken string) (*model.Session, *model.Fail) {
	if accessToken == "" {
		return nil, &model.Fail{
			Message:    "access token cannot be empty",
			StatusCode: http.StatusBadRequest,
		}
	}

	foundSession, err := s.repo.GetByToken(accessToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.Fail{
				Message:    "session not found: " + err.Error(),
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &model.Fail{
			Message:    "failed to get session: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return foundSession, nil
}
