package main

import (
	"net"
	"fmt"
	//	"strings"
	//	"strconv"
)

const	NUM_BYTES = 512

var		server_dead int

var		activeUsers = make(map[int]User)
var		serverMessages = make(chan *Message, 500)

func	main() {
	clientNum := 0
	activeUsers[0] = User{0, "SERVERNAME", "PASSWORD", false, nil, make(chan *Message, 512)}

	// initialize listener
	listener, err := net.Listen("tcp", ":4444")
	if (err != nil) {
		panic(err)
	}
	go serverCloser(listener)
	fmt.Println("Server listening on", listener.Addr())

	go serverDistributeMessages()

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

func	serverDistributeMessages() {
	for ; server_dead == 0 ; {
		select {
		case msg := <-serverMessages:
			fmt.Println("Server recieved message")
			for num, user := range activeUsers {
				//newMsg := *msg
				fmt.Println(num)
				if (msg.sender.IDNumber != user.IDNumber) && (user.IDNumber != 0) && (user.active) {
					user.inbox <-msg //newMsg
				}
			}
		default:
			continue

		}
	}
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
