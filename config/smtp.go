package config

import validation "github.com/go-ozzo/ozzo-validation"

type SMTP struct {
	Available bool   `json:"available" yaml:"available"`
	Host      string `json:"host" yaml:"host"`
	Port      string `json:"port" yaml:"port"`
	From      string `json:"from" yaml:"from"`
	Password  string `json:"password" yaml:"password"`
}

func (config SMTP) Validate() error {
	return validation.ValidateStruct(&config,
		validation.Field(&config.Host, validation.Required),
		validation.Field(&config.Port, validation.Required),
		validation.Field(&config.From, validation.Required),
		validation.Field(&config.Password, validation.Required),
	)
}
