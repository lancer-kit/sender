package config

import validation "github.com/go-ozzo/ozzo-validation"

type Twilio struct {
	Available  bool   `yaml:"available"`
	AccountSid string `yaml:"account_sid"`
	AuthToken  string `yaml:"auth_token"`
	APIURL     string `yaml:"api_url"`
	Sender     string `yaml:"sender"`
}

func (sms Twilio) Validate() error {
	return validation.ValidateStruct(&sms,
		validation.Field(&sms.AccountSid, validation.Required),
		validation.Field(&sms.AuthToken, validation.Required),
		validation.Field(&sms.APIURL, validation.Required),
		validation.Field(&sms.Sender, validation.Required),
	)
}
