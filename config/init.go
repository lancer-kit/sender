package config

import (
	"io/ioutil"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

func Config(path string) (*Cfg, error) {
	rawConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read config file: %s", path)
	}

	cfg := new(Cfg)
	if err = yaml.Unmarshal(rawConfig, cfg); err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshal config file, raw config: %s", rawConfig)
	}

	if err = cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid configs")
	}

	return cfg, nil
}
