package resp

import (
	"sync"
)

type Mux struct {
	rwMutex  sync.RWMutex
	handlers map[string]Handler
}

func NewMux() *Mux {
	return &Mux{handlers: map[string]Handler{}}
}

func (m *Mux) Add(command string, handler Handler) {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	m.handlers[command] = handler
}

func (m *Mux) ServeRESP(req *Req) (RType, error) {
	command := req.Command

	m.rwMutex.RLock()
	handler := m.handlers[command]
	m.rwMutex.RUnlock()

	if handler == nil {
		return nil, errCommandNotSupported(command)
	}
	return handler.ServeRESP(req)
}
