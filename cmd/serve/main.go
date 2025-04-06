package main

import (
	"fmt"
	"io"

	"github.com/quesadelias/http-protocol/internal/http"
	"github.com/quesadelias/http-protocol/internal/request"
)

func main() {
	server := http.New(":8080")

	server.HandleFunc("/", root)

	err := server.Listen()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer server.Close()
}

func root(r request.Request) string {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "HTTP/1.1 500 Internal Server Error\r\nContent-Length: 22\r\n\r\nInternal Server Error"
	}

	fmt.Println(string(body))
	response := "HTTP/1.1 200 OK\r\nContent-Length: 13\r\n\r\nHello, World!"
	return response
}
