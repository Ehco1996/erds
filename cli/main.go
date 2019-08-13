package main

import (
	"log"

	"github.com/Ehco1996/erds"
)

var addr = ":6380"

func main() {
	log.Printf("started server at %s", addr)
	server := erds.NewServer()
	server.ListenAndServe(addr)
}
