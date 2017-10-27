package main

import(
	"net"
	//"fmt"
)

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

type	Message struct {

	Sender		*Client

	Reciever	*Client // should be changed from a single client to a ChatRoom eventually
	// unless we want to keep it, and figure out how to differentiate direct messages and global

	whisper		bool

	Text		*string // does this need to be a pointer?
}

type	Server struct {

	Clients		[]*Client
	// list of all clients ever connected to the server
	// includes not active users

	Connection	chan net.Conn
	// not sure how you're handling this but I'll let you do it and explain it when it works

	Incoming	chan *Message
	// messages sent TO the server to be redistributed as necesssary
	// is this necessary?
}

type	Client struct {

	isActive	bool
	// wether the user is currently connected or not
	// can also be used for "AWAY" functionality maybe?


	nickname	string
	// what is displayed to other users

	username	string
	password	string
	// what is used for authentication

	connection	net.Conn
	// connection to a single user, for sending/recieving data messages

	Incoming	chan *Message
	// data sent from client terminal to the server
	// messages/commands sent BY the user

	Outgoing	chan *Message
	// data to be sent to the client along the channel
	// messages sent TO the user
}
