package main

import (
	"net"
	"fmt"
	"strings"
)

func	newClient(conn net.Conn, serv *Server) (*Client, error) {
	// TODO:
	// handle authentication with the server, checking against previous clients and if the connection just needs to be updated

	// TODO: update this struct with relevant values for nick, user, pass
	client := &Client{
		true,
		"nickname",
		"username",
		"password",
		conn,
		make(chan *Message),
	}
	go setClientInbound(client, serv)

	return client, nil
}

func	setClientInbound(client *Client, serv *Server) {
	err := error(nil)
	buffer := make([]byte, 513)
	var strlen int
	for ; err == nil; {
		strlen, err = client.connection.Read(buffer)
		if (strlen <= 1) {
			err = sendServerResponse("ERROR", "Message length was too short.\n", client.connection)
		} else if (strlen > 512) {
			for {
				// Throw out excess bytes
				if (strings.Index(string(buffer)[0:strlen], "\n") != -1) {
					break
				}
				strlen, err = client.connection.Read(buffer)
			}
			err = sendServerResponse("ERROR", "Message length was too long.\n", client.connection)
		} else if (err != nil) {
			fmt.Println(err)
		} else {
			handleClientInput(client, serv, string(buffer)[0:strlen], strlen)
		}
	}
}

func	handleClientInput(client *Client, serv *Server, input string, strlen int) {
	var prefix, cmd string
	var params []string
	err := error(nil)
	input = strings.TrimSuffix(input, "\r\n")	
	/**  OPTIONAL PREFIX **/
	if (len(input) > 1 && input[0] == ':') {
		fields := strings.SplitN(input, " ", 2)
		if (len(fields) < 2) {
			// Incomplete input
			return
		}
		prefix = fields[0][1:]
		input = strings.TrimLeft(fields[1], " ")
	}
	if (len(input) < 1) {
		// (Remaining) input is empty
		return
	}
	/** COMMAND **/
	fields := strings.SplitN(input, " ", 2)
	if (len(fields[0]) < 1) {
		// Missing or incorrectly spaced command
		return
	}
	// Commands should be uppercase anyway, but just in case...
	cmd = strings.ToUpper(fields[0])
	/** OPTIONAL PARAMETERS **/
	if (len(fields) > 1) {
		input = strings.TrimLeft(fields[1], " ")
		for {
			if (len(input) > 1 && input[0] == ':') {
				// Message
				params = append(params, input[1:])
				break
			}
			// Command parameters
			fields := strings.SplitN(input, " ", 2)
			params = append(params, fields[0])
			if (len(fields) > 1) {
				input = strings.TrimLeft(fields[1], " ")
			} else {
				break
			}
		}
	}
	//fmt.Printf("Prefix: %s Cmd: %s Params: %q\n", prefix, cmd, params)
	msg := makeMessage(client, prefix, cmd, params)
	err = callCommand(msg, serv)
	if (err != nil) {
		fmt.Println(err)
	}
}

func	makeMessage(client *Client, prefix, cmd string, params []string) *Message {
	msg := &Message{client, prefix, cmd, params}
	return (msg)
}

func	sendMessageAlongConnection(prefix, cmd, msg string, conn net.Conn) error {
	var out string
	if (len(prefix) > 0) {
		out += ":"
	}
	out += prefix + " " + cmd + " :" + msg + "\r\n"
	_, err := conn.Write([]byte(out))
	return (err)
}

func	sendServerResponse(cmd, msg string, conn net.Conn) error {
	addr := conn.LocalAddr().String()
	return (sendMessageAlongConnection(addr, cmd, msg, conn))
}

func	sendToRoom(exclude, prefix, cmd, msg string, room *ChatRoom) {
	for _, v := range room.Clients {
		if (v.nickname != exclude) {
			sendMessageAlongConnection(prefix, cmd, msg, v.connection)
		}
	}
}

func	sendServerResponseRoom(exclude, cmd, msg string, room *ChatRoom, conn net.Conn) {
	addr := conn.LocalAddr().String()
	sendToRoom(exclude, addr, cmd, msg, room)
}

func	findRoom(Rooms []*ChatRoom, name string) *ChatRoom {
	var	room *ChatRoom
	for _, v := range Rooms {
		if (v.name == name) {
			room = v;
			break
		}
	}
	return (room)
}

func	findInRoom(room *ChatRoom, nick string) bool {
	inRoom := false
	for _, v := range room.Clients {
		if (v.nickname == nick) {
			inRoom = true
		}
	}
	return (inRoom)
}

func partRoom(clients []*Client, nick string) []*Client {
    for i, v := range clients {
        if v.nickname == nick {
            return append(clients[:i], clients[i+1:]...)
        }
    }
    return clients
}

func	callCommand(msg *Message, serv *Server) error {
	err := error(nil)
	if (msg.cmd == "PRIVMSG") {
		if (len(msg.params) < 2) {
			// Missing target/message
			err = sendServerResponse("ERROR", "Command missing parameters", msg.Sender.connection)
			return (err)
		}
		// TO-DO: Dispatch messages to target
		fmt.Printf("Dest: %s Message: %q\n", msg.params[0], msg.params[1])
	} else if (msg.cmd == "JOIN") {
		if (len(msg.params) < 1 || msg.params[0][0] != '#') {
			err = sendServerResponse("ERROR", "usage: /join #channel\n", msg.Sender.connection)
			return (err)
		}
		room := findRoom(serv.Rooms, msg.params[0])
		if (room != nil) {
			inRoom := findInRoom(room, msg.Sender.nickname)
			if (inRoom) {
				err = sendServerResponse("ERROR", "Already in channel\n", msg.Sender.connection)
				return (err)
			}
		} else {
			room = &ChatRoom{ nil, nil, msg.params[0], ""}
			room.Admins = append(room.Admins, msg.Sender)
			serv.Rooms = append(serv.Rooms, room)
			fmt.Printf("%s created channel %s\n", msg.Sender.nickname, room.name)
		}
		fmt.Printf("%s joined channel %s\n", msg.Sender.nickname, room.name)
		sendServerResponseRoom(msg.Sender.nickname, "PRIVMSG", msg.Sender.nickname + " has joined the channel.", room, msg.Sender.connection)
		room.Clients = append(room.Clients, msg.Sender)
	} else if (msg.cmd == "PART") {
		if (len(msg.params) < 1 || msg.params[0][0] != '#') {
			err = sendServerResponse("ERROR", "usage: /part #channel\n", msg.Sender.connection)
			return (err)
		}
		room := findRoom(serv.Rooms, msg.params[0])
		inRoom := false
		if (room != nil) {
			inRoom = findInRoom(room, msg.Sender.nickname)
		}
		if (room == nil || !inRoom) {
			err = sendServerResponse("ERROR", "Not in channel " + msg.params[0] + "\n", msg.Sender.connection)
			return (err)
		}
		room.Clients = partRoom(room.Clients, msg.Sender.nickname)
		err = sendServerResponse("PRIVMSG", "You left " + room.name + "\n", msg.Sender.connection)
	} else {
		// TO-DO: Dispatch command
		err = sendServerResponse("ERROR", "Unknown command", msg.Sender.connection)
	}
	if (err != nil) {
		fmt.Println(err)
	}
	return (err)
}

