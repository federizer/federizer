package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
	"strconv"
)

const (
	MinPort = 1
	MaxPort = 65535
)

type Config struct {
	ServerHost string `yaml:"serverHost"`
	ServerPort int    `yaml:"serverPort"`
}

/*
	General rule for reading order of configuration values:
	1. Command line.
	2. Config file thats name is declared on the command line.
	3. Environment vars - implemented
	4. Local config file (if exists) - implemented
	5. Global config file (if exists)
*/

func (c *Config) Load(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w; ensure all fields ('host', 'port') are correctly defined", err)
	}

	if host := os.Getenv("SERVER_HOST"); host != "" {
		c.ServerHost = host
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid port value: %s; must be an integer", port)
		}
		c.ServerPort = p
	}

	if c.ServerPort < MinPort || c.ServerPort > MaxPort {
		return fmt.Errorf("invalid port value: %d; port must be between %d and %d", c.ServerPort, MinPort, MaxPort)
	}

	if c.ServerHost == "" {
		return fmt.Errorf("host value is required; expected format: 'hostname' or 'IP address', e.g., 'localhost' or '127.0.0.1'")
	}

	return nil
}
