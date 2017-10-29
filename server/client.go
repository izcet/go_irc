package main

import (
	"net"
	"fmt"
	"strings"
)

func	newClient(conn net.Conn, serv *Server) (*Client, error) {
	// handle authentication with the server, checking against previous clients and if the connection just needs to be updated

	client := &Client{
		true,
		"nickname",
		"username",
		"password",
		conn,
		make(chan *Message),
		make(chan *Message),
	}
	go setClientInbound(client)
	go setClientOutbound(client)

	return client, nil
}

func	setClientInbound(client *Client) {
	err := error(nil)
	buffer := make([]byte, 513)
	var strlen int
	for ; err == nil; {
		strlen, err = client.connection.Read(buffer)
		if (strlen <= 1) {
			err = sendMessageAlongConnection("Error: Message length was too short.\n", client.connection)
		} else if (strlen > 512) {
			for {
				// Throw out excess bytes
				strlen, err = client.connection.Read(buffer)
				if (strings.Index(string(buffer)[0:strlen], "\n") != -1) {
					break
				}
			}
			err = sendMessageAlongConnection("ERROR: Message length was too long.\n", client.connection)
		} else if (err != nil) {
			fmt.Println(err)
		} else {
			handleClientInput(client, string(buffer)[0:strlen], strlen)
		}
	}
}

func	setClientOutbound(client *Client) {
	err := error(nil)
	for ; err == nil ; {
		select {
		case msg := <-client.Outgoing:
			err = sendMessageToClient(msg, client)
		default:
			continue
		}
	}
}

func	handleClientInput(client *Client, input string, strlen int) {
	var prefix, cmd string
	var params []string
	err := error(nil)
/*
	if (input[0] == '/') {
		err = callCommand(client, input, strlen)
	} else {
		msg := makeMessage(client, input)
		client.Incoming <-msg
	}
*/
	/**  OPTIONAL PREFIX **/
	input = strings.TrimSuffix(input, "\r\n")
	if (input[0] == ':') {
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
			if (input[0] == ':') {
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
	fmt.Printf("Prefix: %s Cmd: %s Params: %q\n", prefix, cmd, params)
	if (err != nil) {
		fmt.Println(err)
	}
}

func	makeMessage(client *Client, input string) *Message {
	msg := &Message{client, nil, false, &input}
	return (msg)
}

func	sendMessageToClient(msg *Message, client *Client) error {
	str := "[" + string(msg.Sender.nickname)
	if (msg.whisper) {
		str = str + " whispered to you"
	}
	str = str + "] " + *msg.Text
	return (sendMessageAlongConnection(str, client.connection))
}

func	sendMessageAlongConnection(msg string, conn net.Conn) error {
	_, err := conn.Write([]byte(msg))
	return (err)
}

func	callCommand(client *Client, input string, strlen int) error {
	err := sendMessageAlongConnection("That's a command!\n", client.connection)
	return (err)
}

