package erds

import (
	"log"

	"github.com/tidwall/redcon"
)

func (s *Server) detach(conn redcon.Conn, cmd redcon.Command) {
	detachedConn := conn.Detach()
	log.Printf("connection has been detached")
	go func(c redcon.DetachedConn) {
		defer c.Close()

		c.WriteString("OK")
		c.Flush()
	}(detachedConn)
}

func (s *Server) ping(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("PONG")
}

func (s *Server) quit(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("OK")
	conn.Close()
}

func (s *Server) set(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) != 3 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	s.db.set(string(cmd.Args[1]), cmd.Args[2])
	conn.WriteString("OK")
}

func (s *Server) get(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	val, ok := s.db.get(string(cmd.Args[1]))
	if !ok {
		conn.WriteNull()
	} else {
		conn.WriteBulk(val)
	}
}

func (s *Server) del(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	ok := s.db.delete(string(cmd.Args[1]))
	if !ok {
		conn.WriteInt(0)
	} else {
		conn.WriteInt(1)
	}
}
