package handlers

import (
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
	"github.com/antonrutkevich/simple-memcache/pkg/infrastructure/resp"
	"github.com/sirupsen/logrus"
)

type keyApi struct {
	logger *logrus.Logger
	engine domain.KeyEngine
}

func NewKeyApi(
	logger *logrus.Logger,
	engine domain.KeyEngine,
) *keyApi {
	return &keyApi{logger: logger, engine: engine}
}

func (s *keyApi) Delete() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsMin(req, 1); err != nil {
			return nil, err
		}

		keys := req.Args

		numDeleted := s.engine.Delete(keys)

		return resp.Int(numDeleted), nil
	})
}

func (s *keyApi) Expire() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 2); err != nil {
			return nil, err
		}

		key := req.Args[0]

		seconds, err := parseInt(req, req.Args[1])
		if err != nil {
			return nil, err
		}

		timeoutSet := s.engine.Expire(key, seconds)

		return resp.Int(timeoutSet), nil
	})
}

func (s *keyApi) Ttl() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 1); err != nil {
			return nil, err
		}

		key := req.Args[0]

		return resp.Int(s.engine.Ttl(key)), nil
	})
}
