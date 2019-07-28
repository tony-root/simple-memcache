package log

import (
	"github.com/antonrutkevich/simple-memcache/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

// TODO: consider filtering out invalid config values at earlier stages
func NewLogger(conf config.LogConf) *logrus.Logger {
	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		logrus.Panic(errors.WithMessagef(err, "unknown log level: %s", conf.Level))
	}

	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(level)

	switch conf.Format {
	case "text":
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
			FullTimestamp:   true,
		})
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
		})
	default:
		logrus.Panic(errors.Errorf("unknown log format: %s", conf.Format))
	}

	return log
}
