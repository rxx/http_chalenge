package main

import (
	"errors"
	"fmt"
	"io"
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

	tcp, err := net.Listen("tcp", "0.0.0.0:4221")

	defer tcp.Close()

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := tcp.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		data, err := readData(conn)
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
	fmt.Printf("Request: %s\n", data)

	return blankSuccessResponse()
}

func blankSuccessResponse() []byte {
	return bytes(fmt.Sprintf("%s %s%s%s", PROTOCOL, HTTP_OK, CRLF, CRLF))
}

func bytes(str string) []byte {
	return []byte(str)
}

func readData(conn net.Conn) ([]byte, error) {
	data := make([]byte, 1024)
	size, err := conn.Read(data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("client has disconnected %w", err)
		}
		return nil, fmt.Errorf("error on reading from socket due to %w", err)
	}
	return data[:size], nil
}
