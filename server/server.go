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

var		activeUsers = make(map[int]User)

func	main() {
	clientNum := 0
	listener, err := net.Listen("tcp", ":4444")
	if (err != nil) {
		panic(err)
	}
	server_dead = 0
	fmt.Println("Server listening on", listener.Addr())
	go serverCloser(listener)

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
	defer fmt.Println("client", num, "disconnected)")
	defer conn.Close()

	var buffer = make([]byte, NUM_BYTES)
	
	strlen, err := conn.Write([]byte("Welcome! What is your name?\n"))
	strlen, err = conn.Read(buffer)
	name := string(buffer)[:strlen - 1]
	user := User{num, name, nil, true, conn}
	activeUsers[num] = user
	go closeClient(user) // in case the server shuts down before the client is closed
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
				fmt.Println("Client", num, "shut down the server.")
				conn.Write([]byte("You have shut down the server.\n"))
				break
			} else {
				fmt.Println("Client number", num, "says", str)
				strlen, err = conn.Write(buffer)
			}
		}
	}

	if (err != nil) {
		fmt.Println(err)
	}
}

func	closeClient(user User) {
	for ; ((server_dead == 0) && user.active); {
		continue
	}
	if (user.active) {
		user.conn.Write([]byte("The server has been shut down.\n"))
	}
}

type	myError struct {
	what string
}

func	(e *myError) Error() string {
	return (e.what)
}

func	makeErr(str string) error {
	return (&myError{str})
}
