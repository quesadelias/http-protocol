package http

import (
	"fmt"
	"net"

	"github.com/quesadelias/http-protocol/internal/request"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	routes     map[string]request.Handler
}

func New(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		routes:     make(map[string]request.Handler),
	}
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Println("Server started")

	s.listener = listener

	s.accept()

	return nil
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) HandleFunc(path string, handler request.Handler) {
	fmt.Printf("Register path: %s\n", path)
	s.routes[path] = handler
}

func (s *Server) accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		fmt.Println("Connection accepted: ", conn.RemoteAddr())

		go request.Handle(conn, &s.routes)
	}
}
