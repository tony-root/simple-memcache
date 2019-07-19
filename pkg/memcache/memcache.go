package memcache

import (
	"bufio"
	"fmt"
	"github.com/antonrutkevich/simple-memcache/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

type MemCache struct {
	conf   config.ServerConf
	logger *logrus.Logger
}

func NewMemCache(conf config.ServerConf, logger *logrus.Logger) *MemCache {
	return &MemCache{conf: conf, logger: logger}
}

func (m *MemCache) Run() {
	address := fmt.Sprintf("localhost:%s", m.conf.Port)

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
		go m.handleConn(tcpConn)
	}

}

func (m *MemCache) handleConn(c *net.TCPConn) {
	defer c.Close()

	input := bufio.NewScanner(c)
	for input.Scan() {
		fmt.Fprintln(c, "\t", strings.ToUpper(input.Text()))
	}
}
