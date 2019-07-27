package handlers

import (
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
	"github.com/antonrutkevich/simple-memcache/pkg/infrastructure/resp"
	"github.com/sirupsen/logrus"
)

type hashApi struct {
	logger *logrus.Logger
	engine domain.HashEngine
}

func NewHashApi(
	logger *logrus.Logger,
	engine domain.HashEngine,
) *hashApi {
	return &hashApi{logger: logger, engine: engine}
}

func (s *hashApi) Set() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 3); err != nil {
			return nil, err
		}

		key := req.Args[0]
		field := req.Args[1]
		value := req.Args[2]

		wasNew, err := s.engine.HSet(key, field, value)
		if err != nil {
			return nil, err
		}

		return resp.Int(wasNew), nil
	})
}

func (s *hashApi) MultiSet() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsMin(req, 3); err != nil {
			return nil, err
		}

		key := req.Args[0]

		if err := validateArgsOdd(req); err != nil {
			return nil, err
		}

		entriesList := req.Args[1:]

		entriesMap := make(map[string]string, len(entriesList)/2)

		for i := 0; i < len(entriesList); i += 2 {
			entriesMap[entriesList[i]] = entriesList[i+1]
		}

		err := s.engine.HMultiSet(key, entriesMap)
		if err != nil {
			return nil, err
		}

		return resp.OK(), nil
	})
}

func (s *hashApi) Get() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 2); err != nil {
			return nil, err
		}

		key := req.Args[0]
		field := req.Args[1]

		value, err := s.engine.HGet(key, field)
		if err != nil {
			return nil, err
		}

		if value == "" {
			return resp.Nil(), nil
		}

		return resp.BulkString(value), nil
	})
}

func (s *hashApi) MultiGet() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsMin(req, 2); err != nil {
			return nil, err
		}

		key := req.Args[0]
		fieldsList := req.Args[1:]

		values, err := s.engine.HMultiGet(key, fieldsList)
		if err != nil {
			return nil, err
		}

		return resp.Array(values), nil
	})
}

func (s *hashApi) GetAll() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsExact(req, 1); err != nil {
			return nil, err
		}

		key := req.Args[0]

		values, err := s.engine.HGetAll(key)
		if err != nil {
			return nil, err
		}

		return resp.Array(values), nil
	})
}

func (s *hashApi) Delete() resp.Handler {
	return resp.HandlerFunc(func(req *resp.Req) (resp.RType, error) {
		if err := validateArgsMin(req, 2); err != nil {
			return nil, err
		}

		key := req.Args[0]
		fields := req.Args[1:]

		numDeleted, err := s.engine.HDelete(key, fields)
		if err != nil {
			return nil, err
		}

		return resp.Int(numDeleted), nil
	})
}
