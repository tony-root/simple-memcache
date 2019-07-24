package main

import (
	"github.com/antonrutkevich/simple-memcache/config"
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
	"github.com/antonrutkevich/simple-memcache/pkg/infrastructure/log"
	"github.com/antonrutkevich/simple-memcache/pkg/infrastructure/resp"
	"github.com/antonrutkevich/simple-memcache/pkg/infrastructure/resp/handlers"
	"github.com/pkg/errors"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	conf := config.MustLoad()
	logger := log.CreateLogger(conf.Log)
	engine := domain.NewEngine()

	stringsApi := handlers.NewStringApi(logger, engine)

	mux := resp.NewMux()
	mux.Add("GET", stringsApi.Get())
	mux.Add("SET", stringsApi.Set())

	logger.Infof("Starting memcache %s:%s from %s on port %s\n", version, commit, date, conf.Server.Port)

	server := resp.Server{
		Addr:    conf.Server.Port,
		Logger:  logger,
		Handler: mux,
	}

	err := server.ListenAndServe()

	logger.Fatalf("%+v", errors.Wrap(err, "Internal error"))
}
