package main

import (
	"net"
	//"fmt"
)

func	newClient(conn net.Conn, serv *Server) (*Client, error) {
	// handle authentication with the server, checking against previous clients and if the connection just needs to be updated

	client := &Client{
		true
		"nickname"
		"username"
		"password"
		conn
		make(chan *Message)
		make(chan *Message)
	}
	go setClientInbound(client)
	go setClientOutbound(client)

	return (client, nil)
}

func	setClientInbound(client *Client) {
	err := error(nil)
	buffer := make([]byte, 512)

	for ; err == nil; {
		select {
		case strlen, err = conn.Read(buffer):
			if (strlen > 512) {
				err = sendMessageAlongConnection("ERROR: Message length was too long.\n", client.connection)
			} else if (err != nil) {
				fmt.Println(err)
			} else {
				handleClientInput(client, string(buffer))
			}
		default:
			continue
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

func	handleClientInput(client *Client, input string) {

}

func	sendMessageAlongConnection(msg string, conn net.Conn) error {
	_, err := conn.Write([]byte(msg))
	return (err)
}
