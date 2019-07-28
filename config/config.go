package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type ServerConf struct {
	Port string `validate:"omitempty,numeric"`
}

type LogConf struct {
	Level  string `validate:"omitempty,oneof=error warn info debug"`
	Format string `validate:"omitempty,oneof=text json"`
}

type Config struct {
	Server ServerConf
	Log    LogConf
}

func MustLoad() *Config {
	var config Config

	// TODO: check how to set up mapping rule automatically
	MustBindEnv("server.port", "SERVER_PORT")
	MustBindEnv("log.level", "LOG_LEVEL")
	MustBindEnv("log.format", "LOG_FORMAT")

	if err := viper.Unmarshal(&config); err != nil {
		logrus.Panic(errors.WithMessage(err, "failed to unmarshal config"))
	}

	if err := validator.New().Struct(&config); err != nil {
		logrus.Panic(errors.WithMessage(err, "invalid config"))
	}

	if config.Log.Level == "" {
		config.Log.Level = "info"
	}

	if config.Log.Format == "" {
		config.Log.Format = "text"
	}

	if config.Server.Port == "" {
		config.Server.Port = "9876"
	}

	return &config
}

func MustBindEnv(name string, envName string) {
	if err := viper.BindEnv(name, envName); err != nil {
		logrus.Panic(errors.WithMessage(err, "bind env failed"))
	}
}
