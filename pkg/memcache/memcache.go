package memcache

import (
	"fmt"
	"github.com/antonrutkevich/simple-memcache/config"
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
	"github.com/antonrutkevich/simple-memcache/pkg/infrastructure/transport/telnet"
	"github.com/antonrutkevich/simple-memcache/pkg/usecase"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
)

type MemCache struct {
	conf    config.ServerConf
	logger  *logrus.Logger
	encoder domain.Encoder
	decoder domain.Decoder
	engine  domain.Engine
}

func NewMemCache(
	conf config.ServerConf,
	logger *logrus.Logger,
	encoder domain.Encoder,
	decoder domain.Decoder,
	engine domain.Engine,
) *MemCache {
	return &MemCache{conf: conf, logger: logger, encoder: encoder, decoder: decoder, engine: engine}
}

func (m *MemCache) Run() {
	address := fmt.Sprintf("localhost:%s", m.conf.Port)

	m.engine.Register(usecase.NewStringGet())
	m.engine.Register(usecase.NewStringSet())

	l, err := net.Listen("tcp", address)
	if err != nil {
		m.logger.Panic(errors.WithMessagef(err, "failed to listen at %s", address))
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			m.logger.Info(errors.WithMessagef(err, "failed to accept connection"))
			continue
		}
		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			m.logger.Info(errors.WithMessagef(err, "not a tcp connection"))
		}
		m.handleConn(tcpConn)
	}

}

func (m *MemCache) handleConn(c *net.TCPConn) {
	defer func() {
		err := c.Close()
		if err != nil {
			m.logger.Warn(err, "failed to close conn")
		}
	}()

	for {
		command, err := telnet.ReadCommand(c)
		if err != nil {
			m.logger.Warn(err, "failed to read command, closing")
			_, err := c.Write(m.encoder.EncodeError(err))
			if err != nil {
				m.logger.Warn(err, "failed to write error to connection")
			}
			break
		}

		result := m.engine.Execute(command)

		encoded, err := m.encoder.EncodeResult(result)
		if err != nil {
			_, err := c.Write(m.encoder.EncodeError(err))
			if err != nil {
				m.logger.Warn(err, "failed to encode result")
				break
			}
			continue
		}

		_, err = c.Write(encoded)
		if err != nil {
			m.logger.Warn(err, "failed to write result to connection, closing")
			break
		}
	}
}
