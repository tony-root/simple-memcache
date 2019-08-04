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

		return resp.RInt(numDeleted), nil
	})
}

func (s *keyApi) Expire() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 2); err != nil {
			return nil, err
		}

		key := req.Args[0]

		seconds := parseInt(req, req.Args[1])
		if seconds == nil {
			return nil, errInvalidInteger(req.Command, req.Args[1])
		}

		timeoutSet := s.engine.Expire(key, *seconds)

		return resp.RInt(timeoutSet), nil
	})
}

func (s *keyApi) Ttl() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 1); err != nil {
			return nil, err
		}

		ttl, err := s.engine.Ttl(req.Args[0])
		if err != nil {
			if err == domain.ErrNoTtlForKey {
				return resp.RInt(-1), nil
			}
			if err == domain.ErrTtlKeyNotFound {
				return resp.RInt(-2), nil
			}
		}

		return resp.RInt(ttl), nil
	})
}
