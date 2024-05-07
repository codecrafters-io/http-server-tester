package http_parser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var CRLF = ([]byte)("\r\n")
var SPACE = []byte{' '}
var TAB = []byte{'\t'}
var CR = []byte{'\r'}
var LF = []byte{'\n'}
var NUL = []byte{0}

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

	Headers []Header

	Body []byte
}

func Parse(data []byte) (httpResponse HTTPResponse, readBytesCount int, err error) {
	reader := bytes.NewReader(data)

	response, err := doParseResponse(reader)
	if err != nil {
		return HTTPResponse{}, 0, err
	}

	return response, len(data) - reader.Len(), nil
}

func parseStatusLine(reader *bytes.Reader) (StatusLine, error) {
	var statusLine StatusLine

	offsetBeforeCurrentSection := GetReaderOffset(reader)
	versionB, err := ReadBytes(reader, 9)
	if err == io.EOF {
		return StatusLine{}, IncompleteHTTPResponseError{
			Reader:  reader,
			Message: fmt.Sprintf("Expected: HTTP-version, Received: %q", versionB),
		}
	}
	version := string(versionB)
	// HTTP/1.X
	if version[:8] != "HTTP/1.1" {
		reader.Seek(int64(offsetBeforeCurrentSection), io.SeekStart)
		return StatusLine{}, InvalidHTTPResponseError{
			Reader:  reader,
			Message: fmt.Sprintf("Expected 'HTTP/1.1', Received: %q", version[:8]),
		}
	}
	if version[len(version)-1] != ' ' {
		reader.Seek(int64(offsetBeforeCurrentSection+9), io.SeekStart)
		return StatusLine{}, InvalidHTTPResponseError{
			Reader:  reader,
			Message: "Expected white-space after 'HTTP/1.1'",
		}
	}
	version = version[:8]
	statusLine.Version = version

	offsetBeforeCurrentSection = GetReaderOffset(reader)
	statusB, err := ReadUntilAnyDelimiter(reader, [][]byte{SPACE, TAB})
	if err == io.EOF {
		reader.Seek(int64(-2), io.SeekCurrent)
		return StatusLine{}, IncompleteHTTPResponseError{
			Reader:  reader,
			Message: "Expected reason phrase (example: 'OK' for 200 response) at end of status line",
		}
	}
	statusCode := string(statusB)
	if len(statusCode) != 3 {
		reader.Seek(int64(offsetBeforeCurrentSection), io.SeekStart)
		return StatusLine{}, InvalidHTTPResponseError{
			Reader:  reader,
			Message: "Invalid status-code field length, Expected: 3 digits, Received: " + strconv.Itoa(len(statusCode)),
		}
	}
	intStatusCode, err := strconv.Atoi(statusCode)
	if err != nil {
		reader.Seek(int64(offsetBeforeCurrentSection), io.SeekStart)
		return StatusLine{}, InvalidHTTPResponseError{
			Reader:  reader,
			Message: "Invalid status-code field, Expected integer value, Received: " + statusCode,
		}
	}
	statusLine.StatusCode = intStatusCode

	// Intentionally lax. rfc9112.html#section-4-8
	reasonB, err := ReadUntilCRLF(reader)
	if err == io.EOF {
		return StatusLine{}, IncompleteHTTPResponseError{
			Reader:  reader,
			Message: "Expected CRLF after status line",
		}
	}
	reason := string(reasonB)
	statusLine.Reason = reason
	return statusLine, nil
}

func parseHeaders(reader *bytes.Reader) ([]Header, error) {
	allHeadersFound := false
	headers := []Header{}

	for !allHeadersFound {
		offsetBeforeCRLF := GetReaderOffset(reader)
		possibleHeaderLine, err := ReadUntilCRLF(reader)
		if err == io.EOF {
			if len(possibleHeaderLine) == 0 {
				return []Header{}, IncompleteHTTPResponseError{
					Reader:  reader,
					Message: "Expected CRLF after all headers",
				}
			} else {
				return []Header{}, IncompleteHTTPResponseError{
					Reader:  reader,
					Message: "Expected CRLF after header value",
				}
			}
		}
		if len(possibleHeaderLine) == 0 {
			// Headers finished
			allHeadersFound = true
		} else {
			// We know header is present
			reader.Seek(int64(offsetBeforeCRLF), io.SeekStart)
			key, err := ReadUntil(reader, []byte(":"))
			if err == io.EOF {
				reader.Seek(int64(-2), io.SeekCurrent)
				return []Header{}, IncompleteHTTPResponseError{
					Reader:  reader,
					Message: "Expected ':' after header key",
				}
			}
			if key[len(key)-1] == ' ' {
				// No WhiteSpace before separator
				return []Header{}, InvalidHTTPResponseError{
					Reader:  reader,
					Message: "No whitespace allowed before header separator",
				}
			}
			value, err := ReadUntilCRLF(reader)
			if err == io.EOF {
				return []Header{}, IncompleteHTTPResponseError{
					Reader:  reader,
					Message: "Expected CRLF after header value",
				}
			}
			header := Header{
				// 9110: 5.5-5: Replace CR, LF or NUL with SP
				Key:   string(key),
				Value: strings.TrimSpace(string(ReplaceCharsWithSpace(value, [][]byte{CR, LF, NUL}))),
			}
			headers = append(headers, header)
		}
	}

	return headers, nil
}

func parseContent(reader *bytes.Reader, contentLength int) ([]byte, error) {
	// If content is present
	if contentLength != -1 {
		content, err := ReadBytes(reader, contentLength)
		if err == io.EOF {
			return []byte{}, IncompleteHTTPResponseError{
				Reader:  reader,
				Message: fmt.Sprintf("Expected content of length %d, Received content of length %d", contentLength, len(content)),
			}
		}
		return content, nil
	}
	return []byte{}, nil
}

func doParseResponse(reader *bytes.Reader) (HTTPResponse, error) {
	var httpResponse HTTPResponse

	statusLine, err := parseStatusLine(reader)
	if err != nil {
		return HTTPResponse{}, err
	}
	httpResponse.StatusLine = statusLine

	headers, err := parseHeaders(reader)
	if err != nil {
		return HTTPResponse{}, err
	}
	httpResponse.Headers = headers

	body, err := parseContent(reader, httpResponse.ContentLength())
	if err != nil {
		return HTTPResponse{}, err
	}
	httpResponse.Body = body

	return httpResponse, nil
}

// FindHeader returns a value of a header matching name.
func (response *HTTPResponse) FindHeader(key string) string {
	for _, header := range response.Headers {
		if strings.EqualFold(header.Key, key) {
			return header.Value
		}
	}
	return ""
}

// ContentLength returns the value of the Content-Length header.
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

func (response *HTTPResponse) FormattedString() string {
	var builder strings.Builder

	builder.WriteString(response.StatusLine.Version)
	builder.WriteString(" ")
	builder.WriteString(strconv.Itoa(response.StatusLine.StatusCode))
	builder.WriteString(" ")
	builder.WriteString(response.StatusLine.Reason)
	builder.WriteString("\r\n")

	for _, header := range response.Headers {
		builder.WriteString(header.Key)
		builder.WriteString(": ")
		builder.WriteString(header.Value)
		builder.WriteString("\r\n")
	}

	builder.WriteString("\r\n")
	builder.Write(response.Body)

	return builder.String()
}
