package service

import (
	"fmt"
	"net/http"

	"github.com/MarselBissengaliyev/soilem/configs"
	"github.com/MarselBissengaliyev/soilem/internal/model"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwiloService struct {
	cfg *configs.Config
}

func NewTwiloService(cfg *configs.Config) *TwiloService {
	return &TwiloService{cfg}
}

func (s *TwiloService) SendSMSConfirmation(user *model.User) *model.Fail {
	cfg := s.cfg.Twilio

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSid,
		Password: cfg.AuthToken,
	})

	messageBody := fmt.Sprintf("Your confirmation code is: %s", user.ConfirmationCode)

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(string(user.PhoneNumber))
	params.SetFrom(cfg.FromNumber)
	params.SetBody(messageBody)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return &model.Fail{
			Message:    fmt.Sprintf("failed to send message to number: %s, error: %s", user.PhoneNumber, err.Error()),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}
