package erds

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tidwall/redcon"
)

// Server struct
type Server struct {
	mux *ServeMux
	db  *DB
}

// NewServer init server
func NewServer() *Server {
	// init db and mux
	server := &Server{
		mux: NewServeMux(),
		db:  initDb()}
	// register all handle func
	server.mux.RegisterHandleFunc("ping", server.ping)
	server.mux.RegisterHandleFunc("detach", server.detach)
	server.mux.RegisterHandleFunc("quit", server.quit)
	server.mux.RegisterHandleFunc("set", server.set)
	server.mux.RegisterHandleFunc("get", server.get)
	server.mux.RegisterHandleFunc("del", server.del)

	return server
}

func accept(conn redcon.Conn) bool {
	// this well called on conn accept or deny
	log.Printf("accept: %s", conn.RemoteAddr())
	return true
}

func closed(conn redcon.Conn, err error) {
	// this is called when the connection has been closed
	log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
}

func (s *Server) shutdown() {
	s.mux.aof.saveAofFile()
}

// ListenAndServe start tcp server
func (s *Server) ListenAndServe(addr string) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	go func() {
		err := redcon.ListenAndServe(addr, s.mux.ServeRESP, accept, closed)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	log.Println("Shutdown reds Server save aof...")
	s.shutdown()
}
