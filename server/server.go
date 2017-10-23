package main

import (
	"net"
	"fmt"
	"strings"
	"strconv"
)

const	NUM_BYTES = 512

var		server_dead int

type	User struct {
	IDNumber	int
	nickname	string
	password	string //unused right now
	active		bool
	conn		net.Conn
}

type	Message struct {
	sender		User
	reciever	[]User
	text		string
}

var		activeUsers = make(map[int]User)
var		messages = make(chan Message, 500)

func	main() {
	clientNum := 0
	activeUsers[0] = User{0, "SERVERNAME", "PASSWORD", false, nil}

	// initialize listener
	listener, err := net.Listen("tcp", ":4444")
	if (err != nil) {
		panic(err)
	}
	go serverCloser(listener)
	fmt.Println("Server listening on", listener.Addr())

	//begin main loop
	server_dead = 0
	for ; server_dead == 0; {
		conn, err := listener.Accept()
		if (err != nil) {
			fmt.Println(err)
			continue
		}
		clientNum += 1
		go handleClient(conn, clientNum)
	}

	fmt.Println("Server no longer listening. Shutting down.")
}

func	serverCloser(listener net.Listener) {
	for ; server_dead == 0 ; {
		continue
	}
	listener.Close()
}

func	handleClient(conn net.Conn, num int) { // work on a way to return a value, using channels

	fmt.Println("(client", num, "connected)")
	defer fmt.Println("(client", num, "disconnected)")
	defer conn.Close()

	var buffer = make([]byte, NUM_BYTES)

	strlen, err := conn.Write([]byte("Welcome! What is your name?\n"))
	strlen, err = conn.Read(buffer)
	if (err != nil) {
		fmt.Println("(client", num, "fucked up typing their own name)")
		return
	}
	name := string(buffer)[:strlen - 1]

	user := User{num, name, strconv.Itoa(num), true, conn} // password doesn't matter
	// I literally just set this to cancel the compiler warning about unused strconv

	activeUsers[num] = user
	go clientListen(user)
	makeMessage(0, string(user.nickname + " has joined the server.\n"))

	go deferClientClose(user) // in case the server shuts down before the client is closed
		// insert code to add the client listener to other messages
	// or a channel/broadcast system

	for ; err == nil; {

		//bzero is hard to find
		for x := 0; x < NUM_BYTES; x++ {
			buffer[x] = 0
		}

		//take the message from the client
		strlen, err = conn.Read(buffer)
		if (err != nil) {
			break
		}
		if (strlen != 0) {
			str := string(buffer)
			str = str[:strlen - 1]
			if (strings.Compare(str, "stop") == 0) {
				server_dead = num
				makeMessage(0, string(user.nickname + " shut down the server.\n"))
				break
			} else {
				makeMessage(num, str)
			}
		}
	}

	if (err != nil) {
		fmt.Println(err)
	}
	user.active = false
	makeMessage(0, string(user.nickname + " has disconnected.\n"))
}

func	clientListen(user User) {
	for ; user.active; {
		msg := <-messages
		if (user.active && (msg.sender.IDNumber != user.IDNumber)) {
			var str string
			if (msg.sender.IDNumber == 0) {
				str = msg.text + "\n"
			} else {
				str = "[" + msg.sender.nickname + "] " + msg.text + "\n"
			}
			_, err := user.conn.Write([]byte(str))
			if (err != nil) {
				fmt.Println(err)
			}
		}
	}
}

func	deferClientClose(user User) {
	for ; ((server_dead == 0) && user.active); {
		continue
	}
	if (user.active) {
		user.conn.Write([]byte("The server has been shut down.\n"))
		user.active = false
	}
}

func	makeMessage(user int, text string) {
	fmt.Println("[", user, "]", text)
	msg := Message{activeUsers[user], nil, text}
	messages <- msg
}






////////////////////////////////////////////////////////////////////////////////////////
// here be dragons

type	myError struct {
	what string
}

func	(e *myError) Error() string {
	return (e.what)
}

func	makeErr(str string) error {
	return (&myError{str})
}
