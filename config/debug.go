package config

type Debug struct {
	On                  bool `yaml:"on"`
	EnableLogging       bool `yaml:"enable_logging"`
	DisableRealProvider bool `yaml:"disable_real_provider"`
}
