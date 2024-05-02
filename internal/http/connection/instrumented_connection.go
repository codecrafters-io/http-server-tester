package http_connection

import (
	"fmt"
)

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
		AfterReadResponse: func(value string) {
			fmt.Printf("%sReceived response: %s", logPrefix, value)
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
