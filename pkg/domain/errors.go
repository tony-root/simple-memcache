package domain

import (
	"errors"
	"fmt"
	"github.com/antonrutkevich/simple-memcache/pkg/domain/core"
)

type typeMismatch struct {
	key          string
	expectedType string
}

func errTypeMismatch(key string, expectedType string) *typeMismatch {
	return &typeMismatch{key: key, expectedType: expectedType}
}

func (e *typeMismatch) Error() string {
	return fmt.Sprintf("%s is not of %s type", e.key, e.expectedType)
}

func (e *typeMismatch) ClientErrorCode() core.ClientErrCode {
	return "TYPE_MISMATCH"
}

var (
	ErrTtlKeyNotFound = errors.New("no key found")
	ErrNoTtlForKey    = errors.New("no ttl for key")
)
