package main

import(
	"net"
	"fmt"
)

func main() {
	listener, err := net.Listen("tcp", ":6667")
	server := newServer()
	if err != nil {
		fmt.Printf("cannot listen")
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("cannot accept")
		}
		server.Connection <- conn
	}
}
