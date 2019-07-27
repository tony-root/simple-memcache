package core

type ClientErrCode string

type clientError interface {
	ClientError() (bool, ClientErrCode)
}

func IsClientError(err error) (bool, ClientErrCode) {
	if e, ok := err.(clientError); ok {
		return e.ClientError()
	}
	return false, ""
}
