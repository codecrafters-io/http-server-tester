package http_response

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

func doParseResponse(reader *bytes.Reader) (HTTPResponse, error) {
	var allHeadersFound = false
	var sectionsFound int
	var R HTTPResponse
	var SL StatusLine

	versionB, err := http_utils.ReadUntilAnyDelimiter(reader, [][]byte{http_utils.SPACE, http_utils.TAB})
	if err == io.EOF {
		return HTTPResponse{}, http_utils.IncompleteInputError{
			Reader:  reader,
			Message: "Expected SP between elements of the status line",
		}
	}
	version := string(versionB)
	// HTTP/1.X
	if len(version) != 8 {
		return HTTPResponse{}, http_utils.BadProtocolError{
			Reader:  reader,
			Message: "Invalid HTTP-version field length",
		}
	}
	// ToDo: Assert if version is 1.1 ?
	SL.Version = version
	sectionsFound++

	statusB, err := http_utils.ReadUntilAnyDelimiter(reader, [][]byte{http_utils.SPACE, http_utils.TAB})
	if err == io.EOF {
		return HTTPResponse{}, http_utils.IncompleteInputError{
			Reader:  reader,
			Message: "Expected SP between elements of the status line",
		}
	}
	statusCode := string(statusB)
	if len(statusCode) != 3 {
		return HTTPResponse{}, http_utils.BadProtocolError{
			Reader:  reader,
			Message: "Invalid status-code field length",
		}
	}
	intStatusCode, err := strconv.Atoi(statusCode)
	if err != nil {
		return HTTPResponse{}, http_utils.BadProtocolError{
			Reader:  reader,
			Message: "Invalid status-code field",
		}
	}
	// ToDo: Assert if version is 1.1 ?
	SL.StatusCode = intStatusCode
	sectionsFound++

	// Intentionally lax. rfc9112.html#section-4-8
	reasonB, err := http_utils.ReadUntilCRLF(reader)
	if err == io.EOF {
		return HTTPResponse{}, http_utils.IncompleteInputError{
			Reader:  reader,
			Message: "Expected CRLF after status line",
		}
	}
	reason := string(reasonB)
	SL.Reason = reason
	R.StatusLine = SL
	sectionsFound++

	// XXX: Review with Paul
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
		fmt.Println("Content Length: ", R.ContentLength())
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

// FindHeader returns a value of a header matching name.
func (response *HTTPResponse) FindHeader(key string) string {
	for _, header := range response.Headers {
		if strings.EqualFold(header.Key, key) {
			return header.Value
		}
	}
	return ""
}

// Host returns the value of the Host header
func (response *HTTPResponse) Host() string {
	return response.FindHeader("Host")
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
