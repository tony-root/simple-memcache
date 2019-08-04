package core

import "github.com/pkg/errors"

type ClientErrCode string

type ClientError interface {
	ClientErrorCode() ClientErrCode
}

func GetClientErrorCode(err error) ClientErrCode {
	if e, ok := errors.Cause(err).(ClientError); ok {
		return e.ClientErrorCode()
	}
	return ""
}
