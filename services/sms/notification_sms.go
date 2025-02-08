package sms

import (
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Config struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

type SMSService struct {
	client *twilio.RestClient
	config Config
}

func NewSMSService(config Config) *SMSService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	return &SMSService{
		client: client,
		config: config,
	}
}

func (s *SMSService) SendSMS(to, message string) error {
	params := &twilioApi.CreateMessageParams{
		To:   &to,
		From: &s.config.FromNumber,
		Body: &message,
	}

	_, err := s.client.Api.CreateMessage(params)
	return err
}
