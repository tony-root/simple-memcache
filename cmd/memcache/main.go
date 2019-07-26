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
	logger := log.NewLogger(conf.Log)
	engine := domain.NewEngine()

	stringsApi := handlers.NewStringApi(logger, engine)
	listApi := handlers.NewListApi(logger, engine)
	hashApi := handlers.NewHashApi(logger, engine)
	keyApi := handlers.NewKeyApi(logger, engine)

	mux := resp.NewMux()
	mux.Add("GET", stringsApi.Get())
	mux.Add("SET", stringsApi.Set())

	mux.Add("LPOP", listApi.LeftPop())
	mux.Add("RPOP", listApi.RightPop())
	mux.Add("LPUSH", listApi.LeftPush())
	mux.Add("RPUSH", listApi.RightPush())
	mux.Add("LRANGE", listApi.Range())

	mux.Add("HGET", hashApi.Get())
	mux.Add("HMGET", hashApi.MultiGet())
	mux.Add("HGETALL", hashApi.GetAll())
	mux.Add("HSET", hashApi.Set())
	mux.Add("HMSET", hashApi.MultiSet())
	mux.Add("HDEL", hashApi.Delete())

	mux.Add("DEL", keyApi.Delete())
	mux.Add("EXPIRE", keyApi.Expire())
	mux.Add("TTL", keyApi.Ttl())

	logger.Infof("Starting memcache %s:%s from %s on port %s\n", version, commit, date, conf.Server.Port)

	server := resp.Server{
		Addr:    conf.Server.Port,
		Logger:  logger,
		Handler: mux,
	}

	err := server.ListenAndServe()

	logger.Fatalf("%+v", errors.Wrap(err, "Internal error"))
}
