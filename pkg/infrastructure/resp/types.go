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
type RString string

func (r RString) Marshal() []byte {
	return []byte("+" + string(r) + delimiter)
}

// Bulk String
type RBulkString string

func (r RBulkString) Marshal() []byte {
	byteLen := len(r)
	if byteLen == 0 {
		return nilValue
	}

	return []byte("$" + strconv.Itoa(byteLen) + delimiter + string(r) + delimiter)
}

// Int
type RInt int

func (r RInt) Marshal() []byte {
	return []byte(":" + strconv.Itoa(int(r)) + delimiter)
}

// Array
type RArray []string

func (r RArray) Values() []string {
	return r
}

// TODO: potential hardcode-based optimization for empty/short list
func (r RArray) Marshal() []byte {
	numItems := len(r)
	var builder strings.Builder

	builder.WriteByte('*')
	builder.WriteString(strconv.Itoa(numItems))
	builder.WriteString(delimiter)

	for _, v := range r {
		builder.Write(RBulkString(v).Marshal())
	}

	return []byte(builder.String())
}

// Nil
type rNil struct{}

func RNil() rNil {
	return rNil{}
}

func (rNil) Marshal() []byte {
	return nilValue
}

// Error special case
func MarshalError(errorMessage string) []byte {
	return []byte("-" + errorMessage + delimiter)
}

// Predefined values
const rOK = RString("OK")

func OK() RString {
	return rOK
}
