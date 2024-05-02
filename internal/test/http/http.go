package http

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const OptimalBufferSize = 1500

var CRLF = ([]byte)("\r\n")

type header struct {
	Key   string
	Value string
}

type HTTPParser struct {
	Method, Path, Version string

	Headers      []header
	TotalHeaders int

	contentLength int
}

// Create a new parser
func NewHTTPParser() *HTTPParser {
	return &HTTPParser{
		Headers:       []header{},
		TotalHeaders:  0,
		contentLength: -1,
	}
}

var (
	ErrBadProto    = errors.New("bad protocol")
	ErrMissingData = errors.New("missing data")
)

func (hp *HTTPParser) ParseRequestLine(requestLine []byte) (int, error) {
	var path int
	var version int

	for i := 0; i < len(requestLine); i++ {
		char := requestLine[i]
		// XXX: Other SP required ?
		if char == ' ' || char == '\t' {
			hp.Method = string(requestLine[:i])
			path = i + 1
			break
		}
	}

	// FIXME: Extract to method
	for i := path; i < len(requestLine); i++ {
		char := requestLine[i]
		if char == ' ' || char == '\t' {
			hp.Path = string(requestLine[path:i])
			version = i + 1
			break
		}
	}

	// ToDo: Assert length here or return MissingData
	hp.Version = string(requestLine[version:])

	fmt.Println("Parsed request line: ", hp.Method, hp.Path, hp.Version)

	// XXX: Return total bytes parsed
	return 0, nil
}

func (hp *HTTPParser) ParseHeaderLine(headerLine []byte) (int, error) {
	var key, value string
	var valueIdx int

	for i := 0; i < len(headerLine); i++ {
		char := headerLine[i]
		if char == ':' {
			if headerLine[i-1] == ' ' {
				return 0, ErrBadProto
			}
			key = string(headerLine[:i])
			valueIdx = i + 1
			break
		}
	}

	for i := valueIdx; i < len(headerLine); i++ {
		// 9110: 5.5-5: Replace CR, LF or NUL with SP
		if headerLine[i] == '\r' || headerLine[i] == '\n' || headerLine[i] == 0 {
			headerLine[i] = ' '
		}
	}
	value = string(headerLine[valueIdx:])
	value = strings.TrimSpace(value)

	fmt.Printf("%s:%s\n", key, value)

	hp.Headers = append(hp.Headers, header{Key: key, Value: value})
	hp.TotalHeaders++

	return 0, nil
}

func (hp *HTTPParser) Parse(request []byte) (int, error) {
	var requestLine []byte
	var headerIdx int

	for i := 0; i < len(request); i++ {
		if i+1 < len(request) && bytes.Equal(request[i:i+2], CRLF) {
			requestLine = request[:i]
			headerIdx = i + 2
			break
		}
	}
	fmt.Println("Request Line: ", string(requestLine))
	_, err := hp.ParseRequestLine(requestLine)
	if err != nil {
		return 0, err
	}

	for i := headerIdx; i < len(request); i++ {
		if i+1 < len(request) && bytes.Equal(request[i:i+2], CRLF) {
			header := request[headerIdx:i]
			if len(header) == 0 {
				fmt.Println("End of headers")
				break
			}

			hp.ParseHeaderLine(header)
			headerIdx = i + 2
		}
	}

	fmt.Println("Total Headers: ", hp.TotalHeaders)
	fmt.Println("Headers: ", hp.Headers)

	fmt.Println("Body: ", string(request[headerIdx:]))
	hp.contentLength = hp.ContentLength()

	return 0, nil
}

// Return a value of a header matching name.
func (hp *HTTPParser) FindHeader(key string) string {
	for _, header := range hp.Headers {
		if strings.EqualFold(header.Key, key) {
			return header.Value
		}
	}
	return ""
}

// Return the value of the Host header
func (hp *HTTPParser) Host() string {
	return hp.FindHeader("Host")
}

// Return the value of the Content-Length header.
// A value of -1 indicates the header was not set.
func (hp *HTTPParser) ContentLength() int {
	headerValue := hp.FindHeader("Content-Length")
	if headerValue != "" {
		contentLength, err := strconv.Atoi(headerValue)
		if err != nil {
			return -1
		}
		return contentLength
	}

	return -1
}

func (hp *HTTPParser) Get() bool {
	return strings.EqualFold(hp.Method, "GET")
}

func (hp *HTTPParser) Post() bool {
	return strings.EqualFold(hp.Method, "POST")
}
