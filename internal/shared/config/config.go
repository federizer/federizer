package config

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port string `yaml:"port"`
}

func (c *Config) Load(data []byte) error {
	return yaml.Unmarshal(data, c)
}
