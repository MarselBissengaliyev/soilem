package service

import (
	"errors"
	"net/http"

	"github.com/MarselBissengaliyev/soilem/configs"
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/repo"
	"github.com/jackc/pgx/v5"
)

type EmailCodeService struct {
	cfg  *configs.Config
	repo *repo.Repository
}

func (s *EmailCodeService) SetEmailCode(
	updateEmailCode model.EmailCode,
	userName model.UserName,
) (*model.EmailCode, *model.Fail) {
	if err := updateEmailCode.Validate(); err != nil {
		return nil, &model.Fail{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	emailCode, err := s.repo.EmailCode.SetEmailCode(updateEmailCode, userName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.Fail{
				Message:    err.Error(),
				StatusCode: http.StatusNotFound,
			}
		}

		return nil, &model.Fail{
			Message:    "failed to set email code: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return emailCode, nil
}

func (s *EmailCodeService) SendEmailConfirmation(user *model.User) *model.Fail {
	// smtpServer := 
}
