package http_request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	http_utils "github.com/codecrafters-io/http-server-tester/internal/http/utils"
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

	Headers []Header

	Body []byte
}

func Parse(data []byte) (httpRequest HTTPRequest, readBytesCount int, err error) {
	reader := bytes.NewReader(data)

	request, err := doParseRequest(reader)
	if err != nil {
		return HTTPRequest{}, 0, err
	}

	return request, len(data) - reader.Len(), nil
}

func doParseRequest(reader *bytes.Reader) (HTTPRequest, error) {
	var allHeadersFound = false
	var sectionsFound int
	var R HTTPRequest
	var RL RequestLine

	methodB, err := http_utils.ReadUntilAnyDelimiter(reader, [][]byte{http_utils.SPACE, http_utils.TAB})
	if err == io.EOF {
		return HTTPRequest{}, http_utils.IncompleteInputError{
			Reader:  reader,
			Message: "Expected SP between elements of the status line",
		}
	}
	method := string(methodB)
	if !contains(http_utils.RequestTypes, method) {
		return HTTPRequest{}, http_utils.BadProtocolError{
			Reader:  reader,
			Message: fmt.Sprintf("Invalid method: %s", method),
		}
	}
	RL.Method = method
	sectionsFound++

	pathB, err := http_utils.ReadUntilAnyDelimiter(reader, [][]byte{http_utils.SPACE, http_utils.TAB})
	if err == io.EOF {
		return HTTPRequest{}, http_utils.IncompleteInputError{
			Reader:  reader,
			Message: "Expected SP between elements of the status line",
		}
	}
	path := string(pathB)
	RL.Path = path
	sectionsFound++

	versionB, err := http_utils.ReadUntilCRLF(reader)
	if err == io.EOF {
		return HTTPRequest{}, http_utils.IncompleteInputError{
			Reader:  reader,
			Message: "Expected CRLF after status line",
		}
	}
	version := string(versionB)
	// HTTP/1.X
	if len(version) != 8 {
		return HTTPRequest{}, http_utils.BadProtocolError{
			Reader:  reader,
			Message: "Invalid HTTP-version field length",
		}
	}
	// ToDo: Assert if version is 1.1 ?
	RL.Version = version
	R.RequestLine = RL
	sectionsFound++

	// FIXME: Identical for Request and Response, extract out
	for !allHeadersFound {
		offsetBeforeCRLF := http_utils.GetReaderOffset(reader)
		possibleHeaderLine, err := http_utils.ReadUntilCRLF(reader)
		if err == io.EOF {
			return R, http_utils.IncompleteInputError{
				Reader:  reader,
				Message: "Expected empty line after header section",
			}
		}
		if len(possibleHeaderLine) == 0 {
			// Headers finished
			allHeadersFound = true
			sectionsFound++
		} else {
			// We know header is present
			reader.Seek(int64(offsetBeforeCRLF), io.SeekStart)
			key, err := http_utils.ReadUntil(reader, []byte(":"))
			if err == io.EOF {
				return R, http_utils.IncompleteInputError{
					Reader:  reader,
					Message: "Expected ':' after header key",
				}
			}
			if key[len(key)-1] == ' ' {
				// No WS before separator
				return R, http_utils.BadProtocolError{
					Reader:  reader,
					Message: "No whitespace allowed before header separator",
				}
			}
			value, err := http_utils.ReadUntilCRLF(reader)
			if err == io.EOF {
				return R, http_utils.IncompleteInputError{
					Reader:  reader,
					Message: "Expected CRLF after header value",
				}
			}
			H := Header{
				// 9110: 5.5-5: Replace CR, LF or NUL with SP
				Key:   string(key),
				Value: strings.TrimSpace(string(http_utils.ReplaceCharsWithSpace(value, [][]byte{http_utils.CR, http_utils.LF, http_utils.NUL}))),
			}
			R.Headers = append(R.Headers, H)
		}
	}

	// Content is present
	if R.ContentLength() != -1 {
		content, err := http_utils.ReadBytes(reader, R.ContentLength())
		if err == io.EOF {
			return R, http_utils.IncompleteInputError{
				Reader:  reader,
				Message: "Expected content of length " + strconv.Itoa(R.ContentLength()),
			}
		}
		R.Body = content
	}
	sectionsFound++

	// UnreadDataCheck ?
	if reader.Len() != 0 {
		return R, http_utils.BadProtocolError{
			Reader:  reader,
			Message: "Unexpected data after content",
		}
	}

	// XXX: Required ?
	if sectionsFound != 5 {
		return R, http_utils.BadProtocolError{
			Reader:  reader,
			Message: "Expected 5 sections in response",
		}
	}

	return R, nil
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
