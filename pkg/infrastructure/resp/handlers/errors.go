package handlers

import (
	"fmt"
	"github.com/antonrutkevich/simple-memcache/pkg/domain/core"
)

var argsNumErrCode = core.ClientErrCode("ARGS_NUM")

type noExactArgsNumMatch struct {
	command      string
	expectedArgs int
	gotArgs      int
}

func errNoExactArgsNumMatch(command string, expectedArgs int, gotArgs int) *noExactArgsNumMatch {
	return &noExactArgsNumMatch{command: command, expectedArgs: expectedArgs, gotArgs: gotArgs}
}

func (e *noExactArgsNumMatch) Error() string {
	return fmt.Sprintf("%s requires %d args exactly, got %d", e.command, e.expectedArgs, e.gotArgs)
}

func (e *noExactArgsNumMatch) ClientError() (bool, core.ClientErrCode) {
	return true, argsNumErrCode
}

type notEnoughArgs struct {
	command string
	minArgs int
	gotArgs int
}

func errNotEnoughArgs(command string, minArgs int, gotArgs int) *notEnoughArgs {
	return &notEnoughArgs{command: command, minArgs: minArgs, gotArgs: gotArgs}
}

func (e *notEnoughArgs) Error() string {
	return fmt.Sprintf("%s requires at least %d args, got %d", e.command, e.minArgs, e.gotArgs)
}

func (e *notEnoughArgs) ClientError() (bool, core.ClientErrCode) {
	return true, argsNumErrCode
}

type argsEven struct {
	command string
	gotArgs int
}

func errArgsEven(command string, gotArgs int) *argsEven {
	return &argsEven{command: command, gotArgs: gotArgs}
}

func (e *argsEven) Error() string {
	return fmt.Sprintf("%s requires odd number of args, got %d", e.command, e.gotArgs)
}

func (e *argsEven) ClientError() (bool, core.ClientErrCode) {
	return true, argsNumErrCode
}

type invalidInteger struct {
	command string
	value   string
}

func errInvalidInteger(command string, value string) *invalidInteger {
	return &invalidInteger{command: command, value: value}
}

func (e *invalidInteger) Error() string {
	return fmt.Sprintf("%s expects a valid integer, got '%s'", e.command, e.value)
}

func (e *invalidInteger) ClientError() (bool, core.ClientErrCode) {
	return true, core.ClientErrCode("INVALID_INTEGER")
}
