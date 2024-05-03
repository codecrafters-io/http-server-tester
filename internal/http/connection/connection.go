package http_connection

import (
	"bytes"
	"errors"
	"net"
	"time"

	http_response "github.com/codecrafters-io/http-server-tester/internal/http/parser/response"
)

type HttpConnectionCallbacks struct {
	// BeforeSendRequest is called when a Request is sent.
	// This can be useful for info logs.
	BeforeSendRequest func(string)

	// BeforeSendBytes is called when raw bytes are sent.
	// This can be useful for debug logs.
	BeforeSendBytes func(bytes []byte)

	// AfterBytesReceived is called when raw bytes are read.
	// This can be useful for debug logs.
	AfterBytesReceived func(bytes []byte)

	// AfterReadResponse is called when a Response is decoded from bytes read.
	// This can be useful for success logs.
	AfterReadResponse func(value http_response.HTTPResponse)
}

type HttpConnection struct {
	// Conn is the underlying TCP connection for sending / receiving http messages.
	Conn net.Conn

	// UnreadBuffer contains bytes that have been read but not decoded as a value yet.
	// It can be used to check whether there are any bytes left in the buffer after reading a value.
	UnreadBuffer bytes.Buffer

	// Callbacks is a set of functions that are called at various points in the connection's lifecycle.
	Callbacks HttpConnectionCallbacks
}

func NewHttpConnection(addr string, callbacks HttpConnectionCallbacks) (*HttpConnection, error) {
	conn, err := newConn(addr)

	if err != nil {
		return nil, err
	}

	return &HttpConnection{
		Conn:         conn,
		UnreadBuffer: bytes.Buffer{},
		Callbacks:    callbacks,
	}, nil
}

func (c *HttpConnection) Close() error {
	return c.Conn.Close()
}

func (c *HttpConnection) SendRequest(request []byte) error {
	if c.Callbacks.BeforeSendRequest != nil {
		c.Callbacks.BeforeSendRequest(string(request))
	}

	return c.SendBytes(request)
}

func (c *HttpConnection) SendBytes(bytes []byte) error {
	if c.Callbacks.BeforeSendBytes != nil {
		c.Callbacks.BeforeSendBytes(bytes)
	}

	n, err := c.Conn.Write(bytes)
	if err != nil {
		return err
	}

	if n != len(bytes) {
		return errors.New("failed to write entire bytes to connection")
	}

	return nil
}

func (c *HttpConnection) ReadResponse() (http_response.HTTPResponse, error) {
	return c.ReadResponseWithTimeout(2 * time.Second)
}

func (c *HttpConnection) ReadIntoBuffer() error {
	// We don't want to block indefinitely, so we'll set a read deadline
	c.Conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))

	buf := make([]byte, 1024)
	n, err := c.Conn.Read(buf)

	if n > 0 {
		c.UnreadBuffer.Write(buf[:n])
	}

	return err
}

func (c *HttpConnection) ReadResponseWithTimeout(timeout time.Duration) (http_response.HTTPResponse, error) {
	shouldStopReadingIntoBuffer := func(buf []byte) bool {
		_, _, err := http_response.Parse(buf)

		return err == nil
	}

	c.readIntoBufferUntil(shouldStopReadingIntoBuffer, timeout)

	response, readBytesCount, err := http_response.Parse(c.UnreadBuffer.Bytes())

	if c.Callbacks.AfterBytesReceived != nil && readBytesCount > 0 {
		c.Callbacks.AfterBytesReceived(c.UnreadBuffer.Bytes()[:readBytesCount])
	}

	if err != nil {
		return http_response.HTTPResponse{}, err
	}

	// We've read a response! Let's remove the bytes we've read from the buffer
	c.UnreadBuffer = *bytes.NewBuffer(c.UnreadBuffer.Bytes()[readBytesCount:])

	if c.Callbacks.AfterReadResponse != nil {
		c.Callbacks.AfterReadResponse(response)
	}

	return response, nil
}

func (c *HttpConnection) readIntoBufferUntil(condition func([]byte) bool, timeout time.Duration) {
	deadline := time.Now().Add(timeout)

	for !time.Now().After(deadline) {
		// We'll swallow these errors and try reading again anyway
		_ = c.ReadIntoBuffer()

		if condition(c.UnreadBuffer.Bytes()) {
			break
		} else {
			time.Sleep(10 * time.Millisecond) // Let's wait a bit before trying again
		}
	}
}

func newConn(address string) (net.Conn, error) {
	attempts := 0

	for {
		var err error
		var conn net.Conn

		conn, err = net.Dial("tcp", address)

		if err == nil {
			return conn, nil
		}

		// Already a timeout
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, err
		}

		// 50 * 100ms = 5s
		if attempts > 50 {
			return nil, err
		}

		attempts += 1
		time.Sleep(100 * time.Millisecond)
	}
}
