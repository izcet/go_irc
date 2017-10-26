package main

import (
	"net"
	"fmt"
//	"strings"
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
	defer fmt.Println("(client", num, "disconnected)")//nice logging messages
	defer conn.Close()//nice cleanup

	var buffer = make([]byte, NUM_BYTES)//

	strlen, err := conn.Write([]byte("Welcome! What is your name?\n"))
	strlen, err = conn.Read(buffer)
	if (err != nil) {
		fmt.Println("(client", num, "can't type their own name)")
		return
	}
	name := string(buffer)[:strlen - 1]
	strlen, err = conn.Write([]byte("Welcome [" + name + "], type /help for a list of commands.\n"))
	if (err != nil) {
		fmt.Println("(client", num, "didn't get the welcome message)")
		return
	}

	//set up a user with a name and a pointer to a message struct channel
	user := User{num, name, "", true, conn, make(chan *Message, 512)}

	activeUsers[num] = user //keep track of users in global
	go clientListen(user)
	makeMessage(0, string(user.nickname + " has joined the server"))
	defer ClientClose(user) // in case the server shuts down before the client is closed

	for ; err == nil; {

		//bzero is hard to find
		for i, _ := range buffer {
			buffer[i] = 0
		}

		//take the message from the client
		strlen, err = conn.Read(buffer)
		if (err != nil) {
			break
		}
		if (strlen != 0) && (strlen < 512) {
			str := string(buffer)
			str = str[:strlen - 1]
			if (str == "/stop") {
				server_dead = 1
				makeMessage(0, string(user.nickname + " shut down the server"))
				break
			} else if (str == "/help") {
				sendToClient(conn, "== COMMAND HELP ==\n/help | get help, because you need it\n/kick <name> | kick a person, because they need it\n/list | get a list of online users, so you know who to harass\n/nick <name> | change your name, because your parents weren't clever enough\n/stop | kill the server, because you need to stop\n== NO MORE HELP ==\n")
			} else if (str == "/list") {
				sendToClient(conn, "(a listing of online users)\n")
			} else if (len(str) > 5) && (str[0:6] == "/nick ") {
				sendToClient(conn, "(nickname changed to " + str[6:] + ")\n")
			} else if (len(str) > 5) && (str[0:6] == "/kick ") {
				sendToClient(conn, "(kicked user " + str[6:] + ")\n")			
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
		select {
		case msg := <-user.inbox:
			if (user.active) {
				var str string
				if (msg.sender.IDNumber == 0) {
					str = "<" + msg.text + ">\n"
				} else {
				str = "[" + msg.sender.nickname + "] " + msg.text + "\n"
				}
				user.conn.Write([]byte(str))
			}	
		default:
			continue
		}
	}
}

func	sendToClient(conn net.Conn, str string) {
	conn.Write([]byte(str))
}

func	ClientClose(user User) {
	if (user.active) {
		user.conn.Write([]byte("The server hias been shut down.\n"))
		user.active = false
	}
}

func	makeMessage(user int, text string) {
	fmt.Println("[", user, "] [", text, "]")
	msg := &Message{activeUsers[user], nil, text}
	select {
	case serverMessages <-msg:
		return
		// fmt.Println("Sent message to server")
	default:
		fmt.Println("No message sent")
	}
}

