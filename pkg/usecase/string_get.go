package usecase

import "github.com/antonrutkevich/simple-memcache/pkg/domain"

type commandStringGet struct {
	Key string
}

func (commandStringGet) Name() string {
	return "GET"
}

type stringGet struct{}

func NewStringGet() *stringGet {
	return &stringGet{}
}

func (*stringGet) CommandName() string {
	return "GET"
}

func (*stringGet) Decode(rawCommand *domain.RawCommand) (domain.Command, error) {
	return commandStringGet{Key: "some-key"}, nil
}

func (*stringGet) Process(cache map[string]interface{}, command domain.Command) interface{} {
	commandStringGet, ok := command.(commandStringGet)
	if !ok {
		return domain.ErrorResult{Value: domain.ErrInvalidCommand}
	}

	value := cache[commandStringGet.Key]
	if value == nil {
		return domain.NilResult{}
	}

	strValue, ok := value.(domain.StringEntry)
	if !ok {
		return domain.ErrorResult{Value: domain.ErrKeyTypeMismatch}
	}

	return domain.StringResult{Value: strValue.Value}
}
