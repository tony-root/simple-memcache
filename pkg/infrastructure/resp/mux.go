package resp

import (
	"github.com/antonrutkevich/simple-memcache/pkg/domain"
)

type Mux struct {
	handlers map[string]Handler
}

func NewMux() *Mux {
	return &Mux{handlers: map[string]Handler{}}
}

func (m *Mux) Add(command string, handler Handler) {
	m.handlers[command] = handler
}

func (m *Mux) ServeRESP(req *Req) (RType, error) {
	command := req.Command
	handler := m.handlers[command]
	if handler == nil {
		return nil, domain.Errorf(domain.CodeUnknownCommand, "'%s' command is not supported", command)
	}
	return handler.ServeRESP(req)
}
