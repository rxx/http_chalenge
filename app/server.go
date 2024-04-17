package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
    "log/slog"
    "strings"
)

const (
	PROTOCOL       = "HTTP/1.1"
	HTTP_OK        = "200 OK"
	HTTP_NOT_FOUND = "404 Not Found"
	CRLF           = "\r\n"
)

func main() {
	slog.Info("Start server on port 4221")

	tcp, err := net.Listen("tcp", "0.0.0.0:4221")

	defer tcp.Close()

	if err != nil {
		slog.Error("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := tcp.Accept()
		if err != nil {
			slog.Error("Error accepting connection: ", err)
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
			slog.Error("Error on reading data:", err)
			return
		}
		response := handleRequest(data)

		_, err = conn.Write(response)

		if err != nil {
			slog.Error("Error on writing data: %v\n", err)
			return
		}

	}
}

func handleRequest(data []byte) []byte {
    requestLines := strings.Split(string(data),"\r\n")
    path := strings.Split(requestLines[0], " ")

    if len(path) != 3 {
       return buildResponse("400 BAD REQUEST")
    }
    url := path[1]

    if url == "/" {
        return buildResponse(HTTP_OK)
        }

    return buildResponse(HTTP_NOT_FOUND)

	return buildResponse(HTTP_OK)
}

func buildResponse(status string) []byte {
	return bytes(fmt.Sprintf("%s %s%s%s", PROTOCOL, status, CRLF, CRLF))
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
