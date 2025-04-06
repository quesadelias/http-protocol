package main

import (
	"fmt"
	"io"

	"github.com/quesadelias/http-protocol/internal/http"
	"github.com/quesadelias/http-protocol/internal/request"
)

func main() {
	server := http.New(":8080")

	server.HandleFunc("/", rootHandler)
	server.HandleFunc("/whack", whackHandler)

	err := server.Listen()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer server.Close()
}

func rootHandler(r request.Request) string {
	if r.Method != "POST" {
		response := "HTTP/1.1 405 Method Not Allowed\r\nContent-Length: 18\r\n\r\nMethod Not Allowed"
		return response
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "HTTP/1.1 500 Internal Server Error\r\nContent-Length: 21\r\n\r\nInternal Server Error"
	}

	fmt.Println(string(body))
	response := "HTTP/1.1 200 OK\r\nContent-Length: 13\r\n\r\nHello, World!"
	return response
}

func whackHandler(r request.Request) string {
	if r.Method != "GET" {
		response := "HTTP/1.1 405 Method Not Allowed\r\nContent-Length: 18\r\n\r\nMethod Not Allowed"
		return response
	}

	response := "HTTP/1.1 200 OK\r\nContent-Length: 6\r\n\r\nWhack!"
	return response
}
