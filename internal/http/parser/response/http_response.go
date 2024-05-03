package http_response

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

type StatusLine struct {
	Version    string
	StatusCode int
	Reason     string
}

type HTTPResponse struct {
	StatusLine StatusLine

	Headers      []Header
	TotalHeaders int

	contentLength int

	Content       []byte
	StringContent string
}

var CRLF = ([]byte)("\r\n")

var (
	ErrBadProto    = errors.New("bad protocol")
	ErrMissingData = errors.New("missing data")
)

func parseRequestLine(requestLine []byte) (StatusLine, int, error) {
	var codeIdx int
	var reasonIdx int
	var RL StatusLine

	for i := 0; i < len(requestLine); i++ {
		char := requestLine[i]
		if char == ' ' || char == '\t' {
			version := string(requestLine[:i])
			if len(version) != 8 {
				return StatusLine{}, 0, ErrBadProto
			}
			RL.Version = version
			codeIdx = i + 1
			break
		}
	}

	// FIXME: Extract to method
	// BUG: Convert to int ?
	for i := codeIdx; i < len(requestLine); i++ {
		char := requestLine[i]
		if char == ' ' || char == '\t' {
			statusCode := requestLine[codeIdx:i]
			if len(statusCode) != 3 {
				return StatusLine{}, 0, ErrBadProto
			}
			intStatusCode, err := strconv.Atoi(string(statusCode))
			if err != nil {
				return StatusLine{}, 0, ErrBadProto
			}
			RL.StatusCode = intStatusCode
			reasonIdx = i + 1
			break
		}
	}

	// Intentionally lax. rfc9112.html#section-4-8
	RL.Reason = string(requestLine[reasonIdx:])

	fmt.Println("Parsed status line: ", RL.Version, RL.StatusCode, RL.Reason)

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

func Parse(data []byte) (httpResponse HTTPResponse, readBytesCount int, err error) {
	// reader := bytes.NewReader(data)

	response, _, err := doParseResponse(data)
	if err != nil {
		return HTTPResponse{}, 0, err
	}

	return response, len(data), nil
}

func doParseResponse(request []byte) (HTTPResponse, int, error) {
	var requestLine []byte
	var headerIdx int
	var content []byte
	var statusLineFound, allHeadersFound bool
	var R HTTPResponse

	for i := 0; i < len(request); i++ {
		if i+1 < len(request) && bytes.Equal(request[i:i+2], CRLF) {
			requestLine = request[:i]
			headerIdx = i + 2
			statusLineFound = true
			break
		}
	}

	if !statusLineFound {
		return R, 0, ErrBadProto
	}

	fmt.Println("Status Line: ", string(requestLine))
	RL, _, err := parseRequestLine(requestLine)
	if err != nil {
		return R, 0, err
	}
	R.StatusLine = RL

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
	bodyIdx := headerIdx + 2
	// Content is present
	if R.contentLength != -1 {
		content = request[bodyIdx:]
		fmt.Println("Content Length: ", R.contentLength)
		fmt.Println("Content Length: ", len(content))
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
	R.StringContent = string(content)

	return R, len(request), nil
}

// Return a value of a header matching name.
func (response *HTTPResponse) FindHeader(key string) string {
	for _, header := range response.Headers {
		if strings.EqualFold(header.Key, key) {
			return header.Value
		}
	}
	return ""
}

// Return the value of the Host header
func (response *HTTPResponse) Host() string {
	return response.FindHeader("Host")
}

// Return the value of the Content-Length header.
// A value of -1 indicates the header was not set.
func (response *HTTPResponse) ContentLength() int {
	headerValue := response.FindHeader("Content-Length")
	if headerValue != "" {
		contentLength, err := strconv.Atoi(headerValue)
		if err != nil {
			return -1
		}
		return contentLength
	}

	return -1
}
