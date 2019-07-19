package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type ServerConf struct {
	Port string `validate:"required,numeric"`
}

type LogConf struct {
	Level  string `validate:"required,oneof=error warn info debug"`
	Format string `validate:"required,oneof=text json"`
}

type Config struct {
	Server ServerConf `validate:"required"`
	Log    LogConf    `validate:"required"`
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

	return &config
}

func MustBindEnv(name string, envName string) {
	if err := viper.BindEnv(name, envName); err != nil {
		logrus.Panic(errors.WithMessage(err, "bind env failed"))
	}
}
