package domain

import (
	"fmt"
	"github.com/pkg/errors"
)

type ErrorCode string

var (
	CodeInternal               = ErrorCode("INTERNAL")
	CodeWrongType              = ErrorCode("WRONG_TYPE")
	CodeUnknownCommand         = ErrorCode("UNKNOWN_COMMAND")
	CodeProtocolError          = ErrorCode("PROTOCOL_ERROR")
	CodeWrongNumberOfArguments = ErrorCode("WRONG_NUMBER_OF_ARGS")
	CodeWrongNumber            = ErrorCode("WRONG_NUMBER")
)

type ClientError interface {
	Code() string
	Message() string
}

type memcacheError struct {
	cause   error
	code    ErrorCode
	message string
}

func (e memcacheError) Error() string {
	if e.cause == nil {
		return string(e.code) + " " + e.message
	}

	return errors.WithMessage(e.cause, string(e.code)+" "+e.message).Error()
}

func (e memcacheError) Code() string {
	return string(e.code)
}

func (e memcacheError) Message() string {
	return e.message
}

func WrapError(err error, code ErrorCode, message string) *memcacheError {
	return &memcacheError{cause: err, code: code, message: message}
}

func WrapErrorf(err error, code ErrorCode, message string, args ...interface{}) *memcacheError {
	return &memcacheError{cause: err, code: code, message: fmt.Sprintf(message, args...)}
}

func Error(code ErrorCode, message string) *memcacheError {
	return &memcacheError{code: code, message: message}
}

func Errorf(code ErrorCode, message string, args ...interface{}) *memcacheError {
	return &memcacheError{code: code, message: fmt.Sprintf(message, args...)}
}
