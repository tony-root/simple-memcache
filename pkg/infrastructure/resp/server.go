package resp

import (
	"fmt"
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
)

var respDefaultPort = "9876"

type Server struct {
	Addr    string
	Handler Handler
	Logger  *logrus.Logger
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

// TODO: proper connection handling
func (srv *Server) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			srv.Logger.Info(errors.WithMessagef(err, "failed to accept connection"))
			continue
		}
		srv.handleConn(conn)
	}

}

func (srv *Server) handleConn(c net.Conn) {
	defer func() {
		err := c.Close()
		if err != nil {
			srv.Logger.Warn(err, "failed to close conn")
		}
	}()

	for {
		command, err := ReadCommand(c)
		if err != nil {
			protocolErr := protocolError(err, "failed to read command")
			if ok := writeError(c, protocolErr, srv.Logger); !ok {
				break
			}
			continue
		}

		logger := srv.Logger.WithField("cmd", command).Logger

		result, err := srv.Handler.ServeRESP(NewReq(command))
		if err != nil {
			if ok := writeError(c, err, logger); !ok {
				break
			}
			continue
		}

		marshalled := result.Marshal()
		fmt.Printf("will write %s\n", marshalled)
		_, err = c.Write(marshalled)
		if err != nil {
			logger.WithError(err).Infof("failed to write result to connection, closing")
			break
		}
	}
}

func protocolError(cause error, message string) error {
	return domain.WrapError(cause, domain.CodeProtocolError, message)
}

func writeError(c net.Conn, cause error, logger *logrus.Logger) bool {
	logger.WithError(cause).Infof("writing error")
	_, err := c.Write(MarshalError(cause))
	if err != nil {
		logger.WithError(cause).Warn(errors.WithMessagef(err, "failed to write error"))
		return false
	}
	return true
}

type Req struct {
	Command string
	Args    []string
}

func NewReq(reqArray *rArray) *Req {
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
