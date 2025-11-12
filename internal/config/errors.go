package config

import (
	"errors"
	"fmt"
)

var (
	ErrConfigNotFound       = errors.New("config not found")
	ErrInvalidConfigFormat  = errors.New("invalid config format")
	ErrMissingRequiredField = errors.New("missing required field")
)

type ConfigError struct {
	Field   string
	Message string
	Err     error
}

func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Config error [%s]: %s: %v", e.Field, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("Config error [%s]: %s", e.Field, e.Message)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}
