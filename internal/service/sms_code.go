package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MarselBissengaliyev/soilem/configs"
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/MarselBissengaliyev/soilem/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSCodeService struct {
	cfg  *configs.Config
	repo repo.SMSCode
}

func NewSMSCodeService(cfg *configs.Config, repo repo.SMSCode) *SMSCodeService {
	return &SMSCodeService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *SMSCodeService) SetSMSCode(
	updateSMSCode model.SMSCode,
	userName model.UserName,
) (*model.SMSCode, *model.Fail) {
	if err := updateSMSCode.Validate(); err != nil {
		return nil, &model.Fail{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	smsCode, err := s.repo.SetCode(updateSMSCode, userName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.Fail{
				Message:    err.Error(),
				StatusCode: http.StatusNotFound,
			}
		}

		return nil, &model.Fail{
			Message:    "failed to set sms code: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return smsCode, nil
}

func (s *SMSCodeService) SendSMSConfirmation(to model.UserPhone, code int) *model.Fail {
	cfg := s.cfg.Twilio

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSid,
		Password: cfg.AuthToken,
	})

	messageBody := fmt.Sprintf("Your confirmation code is: %v", code)

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(string(to))
	params.SetFrom(cfg.FromNumber)
	params.SetBody(messageBody)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return &model.Fail{
			Message:    fmt.Sprintf("failed to send message to number: %s, error: %s", to, err.Error()),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}
