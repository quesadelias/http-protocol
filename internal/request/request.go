package request

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Request struct {
	HttpVersion string
	Path        string
	Method      string
	Headers     map[string]string
	Body        io.Reader
}

// TODO: add a response writer: w io.Writer
type Handler func(r Request) string

func Handle(conn net.Conn, routes map[string]Handler) {
	defer conn.Close()

	bufReader := bufio.NewReader(conn)

	requestLine, err := bufReader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request line:", err)
		return
	}

	requestLineParts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(requestLineParts) < 3 {
		fmt.Println("Invalid request line")
		return
	}

	request := Request{
		HttpVersion: requestLineParts[2],
		Path:        requestLineParts[1],
		Method:      requestLineParts[0],
		Headers:     make(map[string]string),
	}

	if _, exists := routes[request.Path]; !exists {
		response := "HTTP/1.1 404 Not Found\r\nContent-Length: 9\r\n\r\nNot Found"
		conn.Write([]byte(response))
		return
	}

	routeHandler := routes[request.Path]

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading headers:", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		headerParts := strings.SplitN(line, ": ", 2)
		if len(headerParts) == 2 {
			request.Headers[headerParts[0]] = headerParts[1]
		}
	}

	if request.Method == "POST" || request.Method == "PUT" {
		contentLengthStr, ok := request.Headers["Content-Length"]
		if !ok {
			fmt.Println("No Content-Length header")
			return
		}

		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			fmt.Println("Invalid Content-Length:", err)
			return
		}

		buffer := make([]byte, 32*1024) // 32KB chunks
		var bytesRead int

		for bytesRead < contentLength {
			toRead := min(len(buffer), contentLength-bytesRead)
			n, err := bufReader.Read(buffer[:toRead])
			if err != nil && err != io.EOF {
				fmt.Println("Error reading body chunk:", err)
				return
			}
			if n == 0 {
				break
			}

			request.Body = bytes.NewReader(buffer[:n])

			bytesRead += n
		}
	}

	response := routeHandler(request)

	// TODO: Handle header and body writes
	conn.Write([]byte(response))
}
