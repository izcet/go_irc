package main

import (
	"net"
)

type	Server struct {

	Clients		[]*Client // list of ALL (including not active) clients ever connected
	Connection	chan net.Conn // @brett, not sure how you're handling this
	Incoming	chan *Message// messages sent TO the server to be redistributed as necesssary
	// is this necessary?
}

type	Client struct {
	isActive	bool
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
	Reciever	*Client // should be changed from a single client to a ChatRoom eventually
	// unless we want to keep it, and figure out how to differentiate direct messages and global
	prefix		string
	cmd			string
	params		[]string
}

type	ChatRoom struct {
	// acts like an IRC or Slack channel
	// not to be confused with channels (the Go data structure)
	Clients		[]*Client
	// all the users currently in the channel
	Admin		*Client
	// the moderator of the channel
	name		string
	description	string
	// pretty self explanatory

	// something else here for properties
	// like invite only, etc.
}
