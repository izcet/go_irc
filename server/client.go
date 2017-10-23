package main

import (
	"net"
	"fmt"
	"strings"
//	"strconv"
)

type	User struct {
	IDNumber	int
	nickname	string
	password	string //unused right now
	active		bool
	conn		net.Conn
	inbox		chan *Message
}

type	Message struct {
	sender		User
	reciever	[]User
	text		string
}

func	handleClient(conn net.Conn, num int) {

	fmt.Println("(client", num, "connected)")
	defer fmt.Println("(client", num, "disconnected)")
	defer conn.Close()

	var buffer = make([]byte, NUM_BYTES)

	strlen, err := conn.Write([]byte("Welcome! What is your name?\n"))
	strlen, err = conn.Read(buffer)
	if (err != nil) {
		fmt.Println("(client", num, "can't type their own name)")
		return
	}
	name := string(buffer)[:strlen - 1]

	user := User{num, name, "", true, conn, make(chan *Message, 512)} // password doesn't matter

	activeUsers[num] = user
	go clientListen(user)
	makeMessage(0, string(user.nickname + " has joined the server"))
	go deferClientClose(user) // in case the server shuts down before the client is closed

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
				makeMessage(0, string(user.nickname + " shut down the server"))
				break
			} else {
				makeMessage(num, str)
			}
		}
	}

	user.active = false
	makeMessage(0, string(user.nickname + " has disconnected"))
}

func	clientListen(user User) {
	for ; user.active; {
		msg := <-user.inbox
		if (user.active) {
			var str string
			if (msg.sender.IDNumber == 0) {
				str = "<" + msg.text + ">\n"
			} else {
				str = "[" + msg.sender.nickname + "] " + msg.text + "\n"
			}
			user.conn.Write([]byte(str))
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
	msg := Message{activeUsers[user], nil, text}
	select {
	case serverMessages <-&msg:
		fmt.Println("Sent message to server")
	default:
		fmt.Println("No message sent")
	}
}

