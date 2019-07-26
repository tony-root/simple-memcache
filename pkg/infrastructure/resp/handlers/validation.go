package handlers

import (
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
	"github.com/antonrutkevich/simple-memcache/pkg/infrastructure/resp"
	"strconv"
)

func validateArgsExact(req *resp.Req, exactArgs int) error {
	if len(req.Args) != exactArgs {
		return domain.Errorf(domain.CodeWrongNumberOfArguments,
			"%s requires %d arguments, got %d", req.Command, exactArgs, len(req.Args))
	}
	return nil
}

func validateArgsMin(req *resp.Req, minArgs int) error {
	if len(req.Args) < minArgs {
		return domain.Errorf(domain.CodeWrongNumberOfArguments,
			"%s requires at least %d arguments, got %d", req.Command, minArgs, len(req.Args))
	}
	return nil
}

func parseInt(stringInt string) (int, error) {
	integer, err := strconv.Atoi(stringInt)
	if err != nil {
		return -1, domain.Errorf(domain.CodeWrongNumber,
			"%s is not an integer or is out of range", stringInt)
	}
	return integer, nil
}
