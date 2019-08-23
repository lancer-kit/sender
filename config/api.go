package config

import validation "github.com/go-ozzo/ozzo-validation"

type API struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	APIRequestTimeout int    `yaml:"api_request_timeout"`
}

func (a API) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Host, validation.Required),
		validation.Field(&a.Port, validation.Required),
		validation.Field(&a.APIRequestTimeout, validation.Required),
	)
}
