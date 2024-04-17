package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
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

	data, err := readData(conn)
	if err != nil {
		slog.Error("Error on reading data:", err)
		return
	}
	response := handleRequest(data)

	_, err = conn.Write(response)

	if err != nil {
		slog.Error("Error on writing data: ", err)
		return
	}
}

func handleRequest(data []byte) []byte {
	slog.Info("Request received: ", "data", string(data))
	requestLines := strings.Split(string(data), "\r\n")
	path := strings.Split(requestLines[0], " ")
	headers := parseHeaders(requestLines[1:])

	if len(path) != 3 {
		return buildBlankResponse("400 BAD REQUEST")
	}
	url := path[1]

	if url == "/" {
		return buildBlankResponse(HTTP_OK)
	} else if strings.HasPrefix(url, "/echo/") {
		echoMsg := strings.Split(url, "/echo/")
		return buildPlainResponse(HTTP_OK, echoMsg[1])
	} else if url == "/user-agent" {
		return buildPlainResponse(HTTP_OK, headers["User-Agent"])
	}

	return buildBlankResponse(HTTP_NOT_FOUND)
}

func parseHeaders(lines []string) map[string]string {
	headers := make(map[string]string)

	for _, line := range lines {
		if strings.Index(line, ":") < 1 {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		headers[parts[0]] = strings.TrimSpace(parts[1])

	}
	return headers
}

func buildBlankResponse(status string) []byte {
	var sb strings.Builder
	sb.WriteString(buildStatusLine(status))
	sb.WriteString(CRLF)

	return bytes(sb.String())
}

func buildPlainResponse(status string, msg string) []byte {
	var sb strings.Builder
	sb.WriteString(buildStatusLine(status))
	sb.WriteString(buildHeader("Content-Type", "text/plain"))
	sb.WriteString(buildHeader("Content-Length", fmt.Sprintf("%d", len(msg))))
	sb.WriteString(CRLF)
	sb.WriteString(msg)
	sb.WriteString(CRLF)

	return bytes(sb.String())
}

func buildStatusLine(status string) string {
	return fmt.Sprintf("%s %s%s", PROTOCOL, status, CRLF)
}

func buildHeader(key, value string) string {
	return fmt.Sprintf("%s: %s%s", key, value, CRLF)
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
