package usecase

import "github.com/antonrutkevich/simple-memcache/pkg/domain"

type commandStringSet struct {
	Key   string
	Value string
}

func (commandStringSet) Name() string {
	return "SET"
}

type stringSet struct{}

func NewStringSet() *stringSet {
	return &stringSet{}
}

func (*stringSet) CommandName() string {
	return "SET"
}

func (*stringSet) Decode(rawCommand *domain.RawCommand) (domain.Command, error) {
	return commandStringGet{Key: "some-key"}, nil
}

func (*stringSet) Process(cache map[string]interface{}, command domain.Command) interface{} {
	stringSetCommand, ok := command.(commandStringSet)
	if !ok {
		return domain.ErrorResult{Value: domain.ErrInvalidCommand}
	}

	cache[stringSetCommand.Key] = domain.StringEntry{Value: stringSetCommand.Value}

	return domain.OkResult
}
