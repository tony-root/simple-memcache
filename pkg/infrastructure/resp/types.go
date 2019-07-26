package resp

import (
	"strconv"
	"strings"
)

var delimiter = "\r\n"
var nilValue = []byte("$-1\r\n")

type RType interface {
	Marshal() []byte
}

// String
type rString struct {
	value string
}

func String(value string) *rString {
	return &rString{value: value}
}

func (r rString) Marshal() []byte {
	return []byte("+" + r.value + delimiter)
}

// Bulk String
type rBulkString struct {
	value string
}

func BulkString(value string) *rBulkString {
	return &rBulkString{value: value}
}

func (r rBulkString) Marshal() []byte {
	byteLen := len(r.value)
	if byteLen == 0 {
		return nilValue
	}

	return []byte("$" + strconv.Itoa(byteLen) + delimiter + r.value + delimiter)
}

// Int
type rInt struct {
	value int
}

func Int(value int) *rInt {
	return &rInt{value: value}
}

func (r rInt) Marshal() []byte {
	return []byte(":" + strconv.Itoa(r.value) + delimiter)
}

// Array
type rArray struct {
	values []string
}

func Array(value []string) *rArray {
	return &rArray{values: value}
}

func (r rArray) Values() []string {
	return r.values
}

// TODO: potential hardcode-based optimization for empty/short list
func (r rArray) Marshal() []byte {
	numItems := len(r.values)
	var builder strings.Builder

	builder.WriteByte('*')
	builder.WriteString(strconv.Itoa(numItems))
	builder.WriteString(delimiter)

	for _, v := range r.values {
		builder.Write(BulkString(v).Marshal())
	}

	return []byte(builder.String())
}

// Nil
type rNil struct{}

func Nil() *rNil {
	return &rNil{}
}

func (rNil) Marshal() []byte {
	return nilValue
}

// Error special case
func MarshalError(error error) []byte {
	value := error.Error()
	return []byte("-" + value + delimiter)
}

// Predefined values
var rOK = rString{value: "OK"}

func OK() *rString {
	return &rOK
}
