package telnet

import (
	"bufio"
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
	"github.com/pkg/errors"
	"io"
	"net"
	"strconv"
)

var (
	ErrInvalidType         = errors.New("invalid type")
	ErrFailedToReadCommand = errors.New("failed to read command")
)

func ReadCommand(c net.Conn) (*domain.RawCommand, error) {
	reader := bufio.NewReader(c)

	arrayType, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if arrayType != '*' {
		return nil, ErrInvalidType
	}

	sizeBytes, err := reader.ReadBytes('\r')
	if err != nil {
		return nil, err
	}

	arraySize, err := strconv.Atoi(string(sizeBytes))
	if err != nil {
		return nil, err
	}

	// Skipping the '\n' from '\r\n'
	_, err = reader.Discard(1)
	if err != nil {
		return nil, err
	}

	var entries = make([]string, 0, arraySize)

	for i := 0; i < arraySize; i++ {
		bulkStringType, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if bulkStringType != '$' {
			return nil, ErrInvalidType
		}

		sizeBytes, err := reader.ReadBytes('\r')
		if err != nil {
			return nil, err
		}

		bytesToRead, err := strconv.Atoi(string(sizeBytes))
		if err != nil {
			return nil, err
		}

		// Skipping the '\n' from '\r\n'
		_, err = reader.Discard(1)
		if err != nil {
			return nil, err
		}

		buf := make([]byte, bytesToRead)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return nil, err
		}

		entries = append(entries, string(buf))

		// Skipping the '\n' from '\r\n'
		_, err = reader.Discard(1)
		if err != nil {
			return nil, err
		}
	}

	if len(entries) == 0 {
		return nil, ErrFailedToReadCommand
	}

	return &domain.RawCommand{
		Name: entries[0],
		Args: entries[1:],
	}, nil
}
