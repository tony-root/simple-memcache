package domain

import "github.com/pkg/errors"

var (
	ErrInternal         = errors.New("internal error")
	ErrUnknownCommand   = errors.New("unknown command")
	ErrMalformedCommand = errors.New("malformed command")
	ErrInvalidCommand   = errors.New("invalid command")
	ErrKeyTypeMismatch  = errors.New("key type mismatch")
)
