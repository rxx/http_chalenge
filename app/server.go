package main

import (
	"fmt"
	"net"
	"os"
)

const (
	PROTOCOL       = "HTTP/1.1"
	HTTP_OK        = "200 OK"
	HTTP_NOT_FOUND = "404 Not Found"
	CRLF           = "\r\n"
)

func main() {
	fmt.Println("Start server on port 4221")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	go handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	var data []byte
	defer conn.Close()

	for {
		_, err := conn.Read(data)
		if err != nil {
			fmt.Printf("Error on reading data: %v\n", err)
			return
		}
		response := handleRequest(data)

		_, err = conn.Write(response)

		if err != nil {
			fmt.Printf("Error on writing data: %v\n", err)
			return
		}

	}
}

func handleRequest(data []byte) []byte {
	var response []byte

	_ = data

	return response
}

func blankSuccessResponse() []byte {
	return bytes(fmt.Sprintf("%s %s%s%s", PROTOCOL, HTTP_OK, CRLF, CRLF))
}

func bytes(str string) []byte {
	return []byte(str)
}
