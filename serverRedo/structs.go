package main

import(
	"net"
	//"fmt"
)

type	Message struct {
	Sender		*Client
	Reciever	*Client
	Text		*string
}

type	Server struct {
	Clients		[]*Client
	Connection	chan net.Conn
	Incoming	chan *Message
	Outgoing	chan *Message
}

type	Client struct {
	nickname	string
	password	string
}
