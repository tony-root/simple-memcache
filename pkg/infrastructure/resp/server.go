package resp

import (
	"github.com/antonrutkevich/simple-memcache/pkg/domain/core"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net"
)

const respDefaultPort = "9876"

type Server struct {
	Addr    string
	Handler Handler
	Logger  *logrus.Logger

	listener net.Listener
	quit     chan bool
	exited   chan bool
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = respDefaultPort
	}
	ln, err := net.Listen("tcp", ":"+addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

func (srv *Server) Serve(l net.Listener) error {
	for {
		// TODO: graceful server shutdown would require a select on done channel here
		c, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				// TODO: are there other errors for which it's fine to skip the connection
				srv.Logger.Info(errors.WithMessagef(err, "failed to accept connection"))
				continue
			}
			return err
		}

		conn := srv.newConn(c)
		go conn.serve(srv.Logger)
	}
}

func (srv *Server) newConn(netConn net.Conn) conn {
	return conn{netConn: netConn, handler: srv.Handler}
}

type conn struct {
	netConn net.Conn
	handler Handler
}

var keyCommand = "cmd"

func (c *conn) serve(logger *logrus.Logger) {
	connection := c.netConn
	addr := connection.RemoteAddr().String()

	logger.Debugf("Connection accepted: %s", addr)
	defer func() {
		logger.Debugf("Connection closed: %s", addr)
		_ = connection.Close()
	}()

	for {
		command, err := ReadCommand(connection)
		if err != nil {
			if isCommonNetReadError(errors.Cause(err)) {
				return
			}
			if ok := c.writeError(err, logger); !ok {
				return
			}
			continue
		}

		commandLogger := logger.WithField(keyCommand, command).Logger

		result, err := c.handler.ServeRESP(NewReq(command))
		if err != nil {
			if ok := c.writeError(err, commandLogger); !ok {
				return
			}
			continue
		}

		_, err = connection.Write(result.Marshal())
		if err != nil {
			c.writeError(err, commandLogger)
			return
		}
	}
}

// Taken from http.Server.
// Checks if the error is a known network error that signals broken connection.
func isCommonNetReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		return true
	}
	if oe, ok := err.(*net.OpError); ok && oe.Op == "read" {
		return true
	}
	return false
}

// Returns false if error is fatal for connection.
func (c *conn) writeError(err error, commandLogger *logrus.Logger) bool {
	if code := core.GetClientErrorCode(err); code != "" {
		if ok := c.writeClientErr(code, err, commandLogger); !ok {
			return false
		}
		return true
	}
	c.writeServerErr(err, commandLogger)
	return false
}

// Returns false if error is fatal for connection.
func (c *conn) writeClientErr(code core.ClientErrCode, cause error, commandLogger *logrus.Logger) bool {
	message := string(code) + " " + cause.Error()
	commandLogger.WithError(cause).Infof("Writing client error")

	_, err := c.netConn.Write(MarshalError(message))
	if err != nil {
		commandLogger.WithError(err).Infof("Failed to write client error")
		return false
	}
	return true
}

// Returns false if error is fatal for connection.
func (c *conn) writeServerErr(cause error, commandLogger *logrus.Logger) {
	message := "ERR_INTERNAL " + cause.Error()
	commandLogger.WithError(cause).Infof("Writing server error")

	_, err := c.netConn.Write(MarshalError(message))
	if err != nil {
		commandLogger.WithError(err).Infof("Failed to write server error")
	}
}

type Req struct {
	Command string
	Args    []string
}

func NewReq(reqArray *RArray) *Req {
	rawArgs := reqArray.Values()
	return &Req{Command: rawArgs[0], Args: rawArgs[1:]}
}

type Handler interface {
	ServeRESP(req *Req) (RType, error)
}

type HandlerFunc func(req *Req) (RType, error)

func (f HandlerFunc) ServeRESP(req *Req) (RType, error) {
	return f(req)
}
