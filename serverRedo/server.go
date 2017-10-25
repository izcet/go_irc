package main

import(
	"net"
	"fmt"
)

type	Message struct {
	Sender *Client
	Reciever *Client
	Text *string
}


type Server struct {
	Clients []*Client
	Connection chan net.Conn
	Incoming chan *Message
	Outgoing chan *Message
}

func (serv *Server) addClient(conn net.Conn) {
	serv.Clients = append(serv.Clients, NewClient(conn))
}

func (serv *Server) sendMessage(msg *Message) {
	for i, user := range serv.Clients {
			user.Incoming <- msg
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
		Connection: make(chan net.Conn)
		Incoming: make(chan *Message)
		Outgoing: make(chan *Message)
	}

	go serv.Listen()

	return serv
}
