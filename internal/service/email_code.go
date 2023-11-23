package service

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"

	"github.com/MarselBissengaliyev/soilem/configs"
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/repo"
	"github.com/jackc/pgx/v5"
	"gopkg.in/gomail.v2"
)

type EmailCodeService struct {
	cfg  *configs.Config
	repo repo.EmailCode
}

func NewEmailCodeService(cfg *configs.Config, repo repo.EmailCode) *EmailCodeService {
	return &EmailCodeService{cfg: cfg, repo: repo}
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

	emailCode, err := s.repo.SetCode(updateEmailCode, userName)
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

func (s *EmailCodeService) SendEmailCode(templatePath string, to string, code int) *model.Fail {
	var body bytes.Buffer

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return &model.Fail{
			Message:    "failed to parse template: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if err = t.Execute(&body, struct{ Code int }{Code: code}); err != nil {
		return &model.Fail{
			Message:    "failed to execute template: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	cfg := s.cfg.Gmail

	m := gomail.NewMessage()

	m.SetHeader("From", cfg.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Email confirmation")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	if err := d.DialAndSend(m); err != nil {
		return &model.Fail{
			Message:    "failed to send confirmation code: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}
