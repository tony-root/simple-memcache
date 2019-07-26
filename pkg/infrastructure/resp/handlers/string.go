package handlers

import (
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
	"github.com/antonrutkevich/simple-memcache/pkg/infrastructure/resp"
	"github.com/sirupsen/logrus"
)

type stringApi struct {
	logger *logrus.Logger
	engine domain.Engine
}

func NewStringApi(
	logger *logrus.Logger,
	engine domain.Engine,
) *stringApi {
	return &stringApi{logger: logger, engine: engine}
}

func (s *stringApi) Set() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 2); err != nil {
			return nil, err
		}

		key := req.Args[0]
		value := req.Args[1]

		s.engine.SetString(key, value)

		return resp.OK(), nil
	})
}

func (s *stringApi) Get() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 1); err != nil {
			return nil, err
		}

		key := req.Args[0]

		result, err := s.engine.GetString(key)
		if err != nil {
			return nil, err
		}

		if result == "" {
			return resp.Nil(), nil
		}

		return resp.BulkString(result), nil
	})
}
