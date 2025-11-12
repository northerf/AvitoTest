package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost    string `mapstructure:"postgres_host"`
	DBUser    string `mapstructure:"postgres_user"`
	DBPass    string `mapstructure:"postgres_password"`
	DBPort    int    `mapstructure:"postgres_port"`
	DBName    string `mapstructure:"postgres_db"`
	SSLMode   string `mapstructure:"postgres_ssl_mode"`
	StartPort string `mapstructure:"start_port"`
}

func (c *Config) Validate() error {
	if c.DBName == "" {
		return fmt.Errorf("%w: DBName", ErrMissingRequiredField)
	}
	if c.DBUser == "" {
		return fmt.Errorf("%w: DBUser", ErrMissingRequiredField)
	}
	if c.DBPass == "" {
		return fmt.Errorf("%w: DBPass", ErrMissingRequiredField)
	}
	if c.StartPort == "" {
		return fmt.Errorf("%w: StartPort", ErrMissingRequiredField)
	}
	return nil
}

func LoadConfig() (Config, error) {
	var config Config
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		var fileNotFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &fileNotFoundErr) {
			return config, fmt.Errorf("%w: %v", ErrConfigNotFound, err)
		}
		return config, fmt.Errorf("failed to read config: %w", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("%w: %v", ErrInvalidConfigFormat, err)
	}
	if err := config.Validate(); err != nil {
		return config, err
	}
	return config, nil
}
