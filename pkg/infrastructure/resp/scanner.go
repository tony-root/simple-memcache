package resp

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"net"
)

var maxArrayLength = 1000
var maxStringLength = 1024 * 1024 // 1mb

func ReadCommand(c net.Conn) (*rArray, error) {
	reader := bufio.NewReader(c)

	arrayLen, err := readArrayLen(reader)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read array length")
	}

	if arrayLen <= 0 {
		return nil, errProtocolError("invalid request args length: %d", arrayLen)
	}

	var entries = make([]string, 0, arrayLen)

	for i := 0; i < arrayLen; i++ {
		bulkString, err := readBulkString(reader)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to read bulk string")
		}

		entries = append(entries, bulkString)
	}

	return Array(entries), nil
}

func readArrayLen(reader *bufio.Reader) (int, error) {
	arrayType, err := reader.ReadByte()
	if err != nil {
		return -1, errors.WithMessage(err, "failed to read array type byte")
	}
	if arrayType != '*' {
		return -1, errProtocolError("want '*' as array identifier, got '%c'", arrayType)
	}

	return readLen(reader, maxArrayLength)
}

func readBulkString(reader *bufio.Reader) (string, error) {
	bulkStringType, err := reader.ReadByte()
	if err != nil {
		return "", errors.WithMessage(err, "failed to read bulk string type byte")
	}
	if bulkStringType != '$' {
		return "", errProtocolError("want '$' as bulk string identifier, got '%c'", bulkStringType)
	}

	bulkStringLength, err := readLen(reader, maxStringLength)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read length")
	}

	buf := make([]byte, bulkStringLength)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return "", errors.WithMessage(err, "failed to read content")
	}

	// Skipping the '\r\n'
	if _, err := reader.Discard(2); err != nil {
		return "", errors.WithMessage(err, "failed to discard line feed")
	}

	return string(buf), nil
}

func readLen(reader *bufio.Reader, maxLength int) (int, error) {
	var length int

	for {
		b, err := reader.ReadByte()
		if b == '\r' {
			break
		}
		if err != nil {
			return -1, errors.WithMessage(err, "failed to read length byte")
		}
		if b < '0' || '9' < b {
			return -1, errProtocolError("want digit byte but got '%c'", b)
		}

		length = 10*length + int(b-'0')

		if length > maxLength {
			return -1, errProtocolError("length exceeds max length of %d", maxLength)
		}
	}

	// Skipping the '\n' of '\r\n'
	if _, err := reader.Discard(1); err != nil {
		return -1, errors.WithMessage(err, "failed to discard line feed")
	}

	return length, nil
}
