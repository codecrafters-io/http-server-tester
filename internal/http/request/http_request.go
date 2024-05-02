package http_request

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Header struct {
	Key   string
	Value string
}

type RequestLine struct {
	Method  string
	Path    string
	Version string
}

type HTTPRequest struct {
	RequestLine RequestLine

	Headers      []Header
	TotalHeaders int

	contentLength int

	Content []byte
}

var CRLF = ([]byte)("\r\n")

var (
	ErrBadProto    = errors.New("bad protocol")
	ErrMissingData = errors.New("missing data")
)

func parseRequestLine(requestLine []byte) (RequestLine, int, error) {
	var pathIndex int
	var versionIndex int
	var RL RequestLine
	// XXX: Update as required
	A := []string{"OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "CONNECT"}

	for i := 0; i < len(requestLine); i++ {
		char := requestLine[i]
		if char == ' ' || char == '\t' {
			method := string(requestLine[:i])
			if method == "" || !contains(A, method) {
				return RequestLine{}, 0, ErrBadProto
			}

			RL.Method = string(requestLine[:i])
			pathIndex = i + 1
			break
		}
	}

	// FIXME: Extract to method
	for i := pathIndex; i < len(requestLine); i++ {
		char := requestLine[i]
		if char == ' ' || char == '\t' {
			path := string(requestLine[pathIndex:i])
			if path == "" {
				return RequestLine{}, 0, ErrBadProto
			}
			RL.Path = path
			versionIndex = i + 1
			break
		}
	}

	// Return detailed error
	// HTTP / DIGIT . DIGIT
	version := (requestLine[versionIndex:])
	if len(version) != 8 {
		return RequestLine{}, 0, ErrBadProto
	}
	RL.Version = string(version)

	fmt.Println("Parsed request line: ", RL.Method, RL.Path, RL.Version)

	return RL, len(requestLine), nil
}

func parseHeaderLine(headerLine []byte) (Header, int, error) {
	var key, value string
	var valueIdx int
	var seperatorFound bool = false
	var H Header

	for i := 0; i < len(headerLine); i++ {
		char := headerLine[i]
		if char == ':' {
			seperatorFound = true
			// No WS before seperator
			if headerLine[i-1] == ' ' {
				return H, 0, ErrBadProto
			}
			key = string(headerLine[:i])
			valueIdx = i + 1
			break
		}
	}
	if !seperatorFound {
		return H, 0, ErrBadProto
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

	H.Key = key
	H.Value = value
	return H, 0, nil
}

func Parse(data []byte) (httpRequest HTTPRequest, readBytesCount int, err error) {
	// reader := bytes.NewReader(data)

	request, _, err := doParseRequest(data)
	if err != nil {
		return HTTPRequest{}, 0, err
	}

	return request, len(data), nil
}

func doParseRequest(request []byte) (HTTPRequest, int, error) {
	var requestLine []byte
	var headerIdx int
	var content []byte
	var requestLineFound, allHeadersFound bool
	var R HTTPRequest

	for i := 0; i < len(request); i++ {
		if i+1 < len(request) && bytes.Equal(request[i:i+2], CRLF) {
			requestLine = request[:i]
			headerIdx = i + 2
			requestLineFound = true
			break
		}
	}

	if !requestLineFound {
		return R, 0, ErrBadProto
	}

	fmt.Println("Request Line: ", string(requestLine))
	RL, _, err := parseRequestLine(requestLine)
	if err != nil {
		return R, 0, err
	}
	R.RequestLine = RL

	for i := headerIdx; i < len(request); i++ {
		if i+1 < len(request) && bytes.Equal(request[i:i+2], CRLF) {
			header := request[headerIdx:i]
			if len(header) == 0 {
				allHeadersFound = true
				break
			}

			H, _, err := parseHeaderLine(header)
			if err != nil {
				return R, 0, err
			}
			R.Headers = append(R.Headers, H)
			R.TotalHeaders++
			// We always point to the next header's starting index
			headerIdx = i + 2
		}
	}

	if !allHeadersFound {
		return R, 0, ErrBadProto
	}

	R.contentLength = R.ContentLength()
	bodyIdx := headerIdx // + 2 ?
	// Content is present
	if R.contentLength != -1 {
		content = request[bodyIdx:]
		if R.contentLength != len(content) {
			return R, 0, ErrMissingData
		}
	} else {
		// No Content-Length header
		content = request[bodyIdx:]
		if len(content) != 0 {
			return R, 0, ErrBadProto
		}
	}

	R.Content = content

	return R, 0, nil
}

// Return a value of a header matching name.
func (hp *HTTPRequest) FindHeader(key string) string {
	for _, header := range hp.Headers {
		if strings.EqualFold(header.Key, key) {
			return header.Value
		}
	}
	return ""
}

// Return the value of the Host header
func (hp *HTTPRequest) Host() string {
	return hp.FindHeader("Host")
}

// Return the value of the Content-Length header.
// A value of -1 indicates the header was not set.
func (hp *HTTPRequest) ContentLength() int {
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

func (hp *HTTPRequest) Get() bool {
	return strings.EqualFold(hp.RequestLine.Method, "GET")
}

func (hp *HTTPRequest) Post() bool {
	return strings.EqualFold(hp.RequestLine.Method, "POST")
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
