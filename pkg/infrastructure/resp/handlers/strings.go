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

func (s stringApi) Set() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		argsRequired := 2

		if len(req.Args) != argsRequired {
			return nil, wrongArgsNumber(req, argsRequired)
		}

		key := req.Args[0]
		value := req.Args[1]

		s.engine.SetString(key, value)

		return resp.OK(), nil
	})
}

func (s stringApi) Get() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		argsRequired := 1

		if len(req.Args) != 1 {
			return nil, wrongArgsNumber(req, argsRequired)
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

func wrongArgsNumber(req *resp.Req, argsRequired int) error {
	return domain.Errorf(domain.CodeWrongNumberOfArguments,
		"%s requires %d arguments, got %d", req.Command, argsRequired, len(req.Args))
}
