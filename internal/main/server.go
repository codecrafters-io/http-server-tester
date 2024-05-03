package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"

	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	http_request "github.com/codecrafters-io/http-server-tester/internal/http/parser/request"
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

	for {
		n, err := c.Read(buf)
		if err != nil {
			return
		}

		fmt.Println("read", n, "bytes")
		fmt.Println(string(buf[0:n]))
		_, _, err = http_request.Parse(buf[0:n])
		if err != nil {
			panic(err)
		}

		c.Write(static)
	}
}

// func main() {
// 	listener, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		panic(err)
// 	}

// 	server := &Server{
// 		Listener: listener,
// 	}
// 	err = server.Listen()
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {
	conn, err := http_connection.NewInstrumentedHttpConnection("localhost:8080", "client")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/", bytes.NewBuffer([]byte("")))
	reqDump, _ := httputil.DumpRequestOut(req, true)
	request, _, err := http_request.Parse(reqDump)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(request)
	conn.SendRequest(reqDump)

	response, err := conn.ReadResponse()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(response)
}
