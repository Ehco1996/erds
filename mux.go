package erds

import (
	"strings"

	"github.com/tidwall/redcon"
)

// HandlerFunc handlefunction
type HandlerFunc func(conn redcon.Conn, cmd redcon.Command)

// ServeMux is an RESP command multiplexer.
type ServeMux struct {
	handlers map[string]HandlerFunc
	aof      *aof
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{
		handlers: make(map[string]HandlerFunc),
		aof:      newAof(),
	}
}

// RegisterHandleFunc registers the handler for the given command.
// If a handler already exists for command, Handle panics.
func (m *ServeMux) RegisterHandleFunc(command string, f HandlerFunc) {
	if command == "" {
		panic("redcon: invalid command")
	}
	if f == nil {
		panic("redcon: nil handler")
	}
	if _, exist := m.handlers[command]; exist {
		panic("redcon: multiple registrations for " + command)
	}

	m.handlers[command] = f
}

// ServeRESP dispatches the command to the handler.
func (m *ServeMux) ServeRESP(conn redcon.Conn, cmd redcon.Command) {
	command := strings.ToLower(string(cmd.Args[0]))
	if f, ok := m.handlers[command]; ok {
		f(conn, cmd)
		if m.aof.isTurnOn() {
			go m.aof.propagateAof(command, cmd)
		}
	} else {
		conn.WriteError("ERR unknown command '" + command + "'")
	}
}
