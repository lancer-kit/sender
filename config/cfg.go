package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/lancer-kit/armory/log"
	"github.com/lancer-kit/armory/natsx"
)

// Cfg main structure of the app configuration.
type Cfg struct {
	API       *API          `yaml:"api"`
	Log       log.Config    `yaml:"log"`
	NATS      *natsx.Config `yaml:"nats"`
	Providers *Providers    `yaml:"providers"`
	Workers   Workers       `yaml:"workers"`
}

func (c Cfg) Validate() (err error) {
	if err = validation.ValidateStruct(&c,
		validation.Field(&c.Log, validation.Required),
		validation.Field(&c.Providers, validation.Required),
		validation.Field(&c.Workers, validation.Required),
	); err != nil {
		return err
	}

	if c.NATS != nil {
		err = c.NATS.Validate()
	}
	if err != nil {
		return err
	}

	if c.API != nil {
		err = c.API.Validate()
	}
	if err != nil {
		return err
	}

	if c.Providers != nil {
		err = c.Providers.Validate()
	}
	if err != nil {
		return err
	}

	return nil
}
