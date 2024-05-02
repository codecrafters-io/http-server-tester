package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/http-server-tester/internal/test/http"
)

type Server struct {
	Listener net.Listener
}

func (s *Server) Listen() error {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		go s.handle(conn)
	}
}

var static = []byte("HTTP/1.1 200 OK\r\nContent-Length: 11\r\n\r\nhello world")

func (s *Server) handle(c net.Conn) {
	fmt.Printf("new conn!\n")
	buf := make([]byte, 1500)

	hp := http.NewHTTPParser()

	for {
		n, err := c.Read(buf)
		if err != nil {
			return
		}

		fmt.Println("read", n, "bytes")
		fmt.Println(string(buf[0:n]))
		_, err = hp.Parse(buf[0:n])
		if err != nil {
			panic(err)
		}

		c.Write(static)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := &Server{
		Listener: listener,
	}
	err = server.Listen()
	if err != nil {
		panic(err)
	}
}
