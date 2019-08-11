package main

import (
	"log"

	"github.com/tidwall/redcon"

	"github.com/Ehco1996/erds"
)

func ping(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("PONG")
}

var addr = ":6380"

func main() {
	log.Printf("started server at %s", addr)

	server := erds.NewServer()
	err := server.ListenAndServe(addr)
	if err != nil {
		log.Fatal(err)
	}
}
