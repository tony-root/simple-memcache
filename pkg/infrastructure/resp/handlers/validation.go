package handlers

import (
	"github.com/antonrutkevich/simple-memcache/pkg/infrastructure/resp"
	"strconv"
)

func validateArgsExact(req *resp.Req, exactArgs int) error {
	if len(req.Args) != exactArgs {
		return errNoExactArgsNumMatch(req.Command, exactArgs, len(req.Args))
	}
	return nil
}

func validateArgsMin(req *resp.Req, minArgs int) error {
	if len(req.Args) < minArgs {
		return errNotEnoughArgs(req.Command, minArgs, len(req.Args))
	}
	return nil
}

func validateArgsOdd(req *resp.Req) error {
	if len(req.Args)%2 == 0 {
		return errArgsEven(req.Command, len(req.Args))
	}
	return nil
}

func parseInt(req *resp.Req, stringInt string) *int {
	integer, err := strconv.Atoi(stringInt)
	if err != nil {
		return nil
	}
	return &integer
}
