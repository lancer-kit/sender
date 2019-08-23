package config

import validation "github.com/go-ozzo/ozzo-validation"

type Mailgun struct {
	Available  bool   `json:"available" yaml:"available"`
	Domain     string `json:"domain" yaml:"domain"`
	Sender     string `json:"sender" yaml:"sender"`
	PrivateKey string `json:"privateKey" yaml:"private_key"`
	PublicKey  string `json:"publicKey" yaml:"public_key"`
}

func (config Mailgun) Validate() error {
	return validation.ValidateStruct(&config,
		validation.Field(&config.Domain, validation.Required),
		validation.Field(&config.Sender, validation.Required),
		validation.Field(&config.PrivateKey, validation.Required),
		validation.Field(&config.PublicKey, validation.Required),
	)
}
