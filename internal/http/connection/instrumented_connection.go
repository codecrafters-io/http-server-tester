package http_connection

import (
	"fmt"

	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
)

// ToDo: logPrefix is required ?
func defaultCallbacks(logPrefix string) HttpConnectionCallbacks {
	return HttpConnectionCallbacks{
		BeforeSendRequest: func(command string) {
			fmt.Printf("%sSent request: %s", logPrefix, command)
		},
		BeforeSendBytes: func(bytes []byte) {
			fmt.Printf("%sSent bytes: %q", logPrefix, string(bytes))
		},
		AfterBytesReceived: func(bytes []byte) {
			fmt.Printf("%sReceived bytes: %q", logPrefix, string(bytes))
		},
		AfterReadResponse: func(value http_parser.HTTPResponse) {
			fmt.Printf("%sReceived response: %v", logPrefix, value)
		},
	}
}

func NewInstrumentedHttpConnection(addr string, clientIdentifier string) (*HttpConnection, error) {
	logPrefix := ""
	if clientIdentifier != "" {
		logPrefix = clientIdentifier + ": "
	}
	return NewHttpConnection(
		addr, defaultCallbacks(logPrefix),
	)
}
