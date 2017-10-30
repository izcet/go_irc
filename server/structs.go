package main

import (
	"net"
)

type	Server struct {
	Clients		[]*Client
	// list of all clients ever connected to the server
	// includes not active users

	Connection	chan net.Conn
	// not sure how you're handling this but I'll let you do it and explain it when it works
	Rooms		[]*ChatRoom
}

type	Client struct {
	isActive	bool
	// wether the user is currently connected or not
	// can also be used for "AWAY" functionality maybe?

	nickname	string
	username	string
	password	string
	connection	net.Conn
	Incoming	chan *Message
	// data sent from client terminal to the server
	// messages/commands sent BY the user
	Outgoing	chan *Message
	// data to be sent to the client along the channel
	// messages sent TO the user
}

type	Message struct {
	Sender		*Client
	prefix		string
	cmd			string
	params		[]string
}

type	ChatRoom struct {
	// acts like an IRC or Slack channel
	// not to be confused with channels (the Go data structure)
	Clients		[]*Client
	// all the users currently in the channel
	Admins		[]*Client
	// the moderator of the channel
	name		string
	description	string
	// pretty self explanatory

	// something else here for properties
	// like invite only, etc.
}
