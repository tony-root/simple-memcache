package telnet

import (
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
	"strconv"
	"strings"
)

var delimiter = "\r\n"
var nilValue = []byte("$-1\r\n")

type encoder struct{}

func NewTelnetEncoder() *encoder {
	return &encoder{}
}

func (e *encoder) EncodeResult(result interface{}) ([]byte, error) {
	var encoded []byte = nil

	switch actual := result.(type) {
	case domain.StringResult:
		encoded = e.EncodeString(actual.Value)
	case domain.BulkStringResult:
		encoded = e.EncodeBulkString(actual.Value)
	case domain.ArrayResult:
		encoded = e.EncodeArray(actual.Value)
	case domain.IntResult:
		encoded = e.EncodeInt(actual.Value)
	case domain.NilResult:
		encoded = e.EncodeNil()
	default:
		return nil, domain.ErrInternal
	}

	return encoded, nil
}

// TODO: potential hardcode-based optimization for empty/short list
func (*encoder) EncodeArray(value []string) []byte {
	numItems := len(value)
	header := "*" + strconv.Itoa(numItems)

	var builder strings.Builder
	for _, v := range value {
		builder.WriteString(v + delimiter)
	}

	return []byte(header + builder.String())
}

func (*encoder) EncodeInt(value int) []byte {
	return []byte(":" + strconv.Itoa(value) + delimiter)
}

func (*encoder) EncodeBulkString(value string) []byte {
	return []byte("$" + value + delimiter)
}

func (*encoder) EncodeString(value string) []byte {
	return []byte("+" + value + delimiter)
}

func (*encoder) EncodeError(error error) []byte {
	value := error.Error()
	return []byte("-" + value + delimiter)
}

func (*encoder) EncodeNil() []byte {
	return nilValue
}

func NewEncoder() *encoder {
	return &encoder{}
}
