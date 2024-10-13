package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"gopkg.in/yaml.v2"
)

const (
	MinPort = 1
	MaxPort = 65535
)

type Config struct {
	ServerHost string `yaml:"serverHost"`
	ServerPort string `yaml:"serverPort"`
}

// toSnakeCase converts a string to snake_case.
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
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

	for i := 0; i < reflect.TypeOf(*c).NumField(); i++ {
		field := reflect.TypeOf(*c).Field(i)
		if value, ok := field.Tag.Lookup("yaml"); ok {
			if len(os.Getenv(strings.ToUpper(toSnakeCase(value)))) > 0 {
				reflect.ValueOf(c).Elem().FieldByName(field.Name).Set(reflect.ValueOf(os.Getenv(strings.ToUpper(toSnakeCase(value)))))
			}
		}
	}

	port, err := strconv.Atoi(c.ServerPort)
	if err != nil {
		return fmt.Errorf("invalid port value: %s; must be an integer", c.ServerPort)
	}

	if port < MinPort || port > MaxPort {
		return fmt.Errorf("invalid port value: %d; port must be between %d and %d", port, MinPort, MaxPort)
	}

	if c.ServerHost == "" {
		return fmt.Errorf("host value is required; expected format: 'hostname' or 'IP address', e.g., 'localhost' or '127.0.0.1'")
	}

	return nil
}
