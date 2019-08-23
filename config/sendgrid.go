package config

import validation "github.com/go-ozzo/ozzo-validation"

type Sendgrid struct {
	Available  bool   `json:"available" yaml:"available"`
	Sender     string `json:"sender" yaml:"sender"`
	SenderName string `json:"senderName" yaml:"sender_name"`
	PrivateKey string `json:"privateKey" yaml:"private_key"`
	Endpoint   string `json:"endpoint" yaml:"endpoint"`
	Host       string `json:"host" yaml:"host"`
}

func (c Sendgrid) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Sender, validation.Required),
		validation.Field(&c.SenderName, validation.Required),
		validation.Field(&c.PrivateKey, validation.Required),
		validation.Field(&c.Endpoint, validation.Required),
		validation.Field(&c.Host, validation.Required),
	)
}
