package telnet

import "github.com/antonrutkevich/simple-memcache/pkg/domain"

type decoder struct{}

func NewDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(rawCommand []byte) (*domain.RawCommand, error) {
	return &domain.RawCommand{
		Name: "GET",
		Args: []string{"key1"},
	}, nil
}
