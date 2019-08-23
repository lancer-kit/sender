package config

import validation "github.com/go-ozzo/ozzo-validation"

type Viber struct {
	Available bool   `yaml:"available"`
	APIURL    string `yaml:"api_url"`
}

func (v Viber) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.APIURL, validation.Required),
	)
}
