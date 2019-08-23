package config

import validation "github.com/go-ozzo/ozzo-validation"

type Whatsapp struct {
	Available bool   `yaml:"available"`
	APIURL    string `yaml:"api_url"`
	APIKey    string `yaml:"api_key"`
}

func (config Whatsapp) Validate() error {
	return validation.ValidateStruct(&config,
		validation.Field(&config.APIURL, validation.Required),
		validation.Field(&config.APIKey, validation.Required),
	)
}
