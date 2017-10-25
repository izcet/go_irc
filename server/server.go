package main

import (
	"net"
	"fmt"
	//	"strings"
	//	"strconv"
)
//Global state bad
const	NUM_BYTES = 512 //maximum size of message

var		server_dead int

var		activeUsers = make(map[int]User) //global map for something with users???
var		serverMessages = make(chan *Message, 500) //global buffered *Message channel ???

func	main() {
	clientNum := 0 //why keep track of total clients
	activeUsers[0] = User{0, "SERVERNAME", "PASSWORD", false, nil, make(chan *Message, 512)} //the first member of the user map is the server possibly for server connection

	// initialize listener
	listener, err := net.Listen("tcp", ":4444") //how are you gonna have a server without a listener 10/10
	if (err != nil) { //maybe the port is taked already
		panic(err)
	}
	defer listener.Close()//close the listener after it's use
	fmt.Println("Server listening on", listener.Addr())//helps to know the address to tell other people

	go serverDistributeMessages() //the messages need to find their way around the network

	//begin main loop
	server_dead = 0//why ?????
	for ; server_dead == 0; {
		conn, err := listener.Accept()
		if (err != nil) {
			fmt.Println(err)
			continue
		}
		clientNum += 1
		go handleClient(conn, clientNum)//need concurrent routines for each client
	}

	fmt.Println("Server no longer listening. Shutting down.")
}

func	serverDistributeMessages() {
	for ; server_dead == 0 ; {//main loop again
		select {
		case msg := <-serverMessages: //if a message appears it has to find it's way to everyone it applies to
			fmt.Println("Server recieved message")
			for num, user := range activeUsers {
				//newMsg := *msg
				fmt.Println(num)//
				if (msg.sender.IDNumber != user.IDNumber) && (user.IDNumber != 0) && (user.active) {// user cant send themselves a message otherwise put it in everyones mailbox
					user.inbox <-msg //newMsg
				}
			}
		default://the serverMessages isn't always busy
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
