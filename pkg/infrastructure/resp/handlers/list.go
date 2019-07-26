package handlers

import (
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
	"github.com/antonrutkevich/simple-memcache/pkg/infrastructure/resp"
	"github.com/sirupsen/logrus"
)

type listApi struct {
	logger *logrus.Logger
	engine domain.ListEngine
}

func NewListApi(
	logger *logrus.Logger,
	engine domain.ListEngine,
) *listApi {
	return &listApi{logger: logger, engine: engine}
}

func (s *listApi) LeftPop() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 1); err != nil {
			return nil, err
		}

		key := req.Args[0]

		result, err := s.engine.LLeftPop(key)
		if err != nil {
			return nil, err
		}

		if result == "" {
			return resp.Nil(), nil
		}

		return resp.BulkString(result), nil
	})
}

func (s *listApi) RightPop() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 1); err != nil {
			return nil, err
		}

		key := req.Args[0]

		result, err := s.engine.LRightPop(key)
		if err != nil {
			return nil, err
		}

		if result == "" {
			return resp.Nil(), nil
		}

		return resp.BulkString(result), nil
	})
}

func (s *listApi) LeftPush() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsMin(req, 2); err != nil {
			return nil, err
		}

		key := req.Args[0]
		values := req.Args[1:]

		listSize, err := s.engine.LLeftPush(key, values)
		if err != nil {
			return nil, err
		}

		return resp.Int(listSize), nil
	})
}

func (s *listApi) RightPush() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsMin(req, 2); err != nil {
			return nil, err
		}

		key := req.Args[0]
		values := req.Args[1:]

		listSize, err := s.engine.LRightPush(key, values)
		if err != nil {
			return nil, err
		}

		return resp.Int(listSize), nil
	})
}

func (s *listApi) Range() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 3); err != nil {
			return nil, err
		}

		key := req.Args[0]

		min, err := parseInt(req.Args[1])
		if err != nil {
			return nil, err
		}

		max, err := parseInt(req.Args[2])
		if err != nil {
			return nil, err
		}

		values, err := s.engine.LRange(key, min, max)
		if err != nil {
			return nil, err
		}

		return resp.Array(values), nil
	})
}
