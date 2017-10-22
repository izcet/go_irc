package main

import (
	"net"
	"fmt"
	"strings"
	"strconv"
)

const	NUM_BYTES = 512

var		server_dead int

func	main() {
	clientNum := 0
	listener, err := net.Listen("tcp", ":4444")
	if (err != nil) {
		panic(err)
	}
	defer listener.Close()

	server_dead = 0
	fmt.Println("Server listening on", listener.Addr())
	for ; server_dead == 0; {
		conn, err := listener.Accept()
		if (err != nil) {
			fmt.Println(err)
			continue
		}
		if (server_dead != 0) {
			_,err := conn.Write([]byte("This server is already closed.\n"))
			if (err != nil) {
				fmt.Println(err)
			}
			conn.Close()
			break
		}
		clientNum += 1
		fmt.Println("Client", clientNum, "connected!")
		go handleClient(conn, clientNum)
	}
}

func	handleClient(conn net.Conn, num int) { // work on a way to return a value, using channels
	defer conn.Close()
	strlen, err := conn.Write([]byte("Welcome! You are client number " + strconv.Itoa(num) + ".\n"))
	var buffer = make([]byte, NUM_BYTES)
	for ; err == nil; {
		for x := 0; x < NUM_BYTES; x++ {
			buffer[x] = 0
		}
		strlen, err = conn.Read(buffer)
		if (server_dead != 0) {
			msg := "Client " + strconv.Itoa(server_dead) + " shut down the server!\n"
			conn.Write([]byte(msg))
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
	fmt.Println("(client", num, "disconnected)")
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
