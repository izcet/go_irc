package main

import(
	"net"
	"fmt"
)

func (serv *Server) addClient(conn net.Conn) {
	client, err := newClient(conn, serv)
	if (err != nil) {
		// just in case creating a client causes an error
		// client should be nil
		fmt.Println(err)
	}
	if (client != nil) {
		serv.Clients = append(serv.Clients, client)
	}
}

func (serv *Server) sendMessage(msg *Message) {
	for _, user := range serv.Clients {
			user.Outgoing <- msg
	}
}

func (serv *Server) Listen() {
	for {
		select {
		case conn := <-serv.Connection:
			serv.addClient(conn)
		case msg := <-serv.Incoming:
			serv.sendMessage(msg)
		default:
			continue
		}
	}
}

func newServer() *Server {
	serv := &Server{
		Connection: make(chan net.Conn),
		Incoming: make(chan *Message),
//		Outgoing: make(chan *Message),
	}

	go serv.Listen()

	return serv
}
