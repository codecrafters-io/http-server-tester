package http_parser

import (
	"bytes"
	"fmt"
	"io"
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

	offsetBeforeCurrentSection := GetReaderOffset(reader)
	versionB, err := ReadUntilAnyDelimiter(reader, [][]byte{SPACE, TAB})
	if err == io.EOF {
		return HTTPResponse{}, IncompleteInputError{
			Reader:  reader,
			Message: fmt.Sprintf("Expected: HTTP-version, Received: %q", versionB),
		}
	}
	version := string(versionB)
	// HTTP/1.X
	if len(version) != 8 {
		// Seek to ending point of last sub-section
		reader.Seek(int64(-2), io.SeekCurrent)
		return HTTPResponse{}, BadProtocolError{
			Reader:  reader,
			Message: "Invalid HTTP-version field length",
		}
	}
	// HTTP/
	if version[0:5] != "HTTP/" {
		// Seek to starting position for current sub-section
		reader.Seek(int64(offsetBeforeCurrentSection+4), io.SeekStart)
		return HTTPResponse{}, BadProtocolError{
			Reader:  reader,
			Message: "Expected HTTP-version field to start with 'HTTP/'",
		}
	}
	// 1.X
	if version[5] != '1' || version[6] != '.' {
		reader.Seek(int64(offsetBeforeCurrentSection+5), io.SeekStart)
		return HTTPResponse{}, BadProtocolError{
			Reader:  reader,
			Message: "Expected HTTP-version 1.X, Received: " + version[5:],
		}
	}

	SL.Version = version
	sectionsFound++

	// offsetBeforeCurrentSection = GetReaderOffset(reader)
	statusB, err := ReadUntilAnyDelimiter(reader, [][]byte{SPACE, TAB})
	if err == io.EOF {
		return HTTPResponse{}, IncompleteInputError{
			Reader:  reader,
			Message: "Status line has missing sections, Expected: HTTP-version status-code reason-phrase",
		}
	}
	statusCode := string(statusB)
	if len(statusCode) != 3 {
		reader.Seek(int64(-2), io.SeekCurrent)
		return HTTPResponse{}, BadProtocolError{
			Reader:  reader,
			Message: "Invalid status-code field length, Expected: 3 digits, Received: " + strconv.Itoa(len(statusCode)),
		}
	}
	intStatusCode, err := strconv.Atoi(statusCode)
	if err != nil {
		reader.Seek(int64(-2), io.SeekCurrent)
		return HTTPResponse{}, BadProtocolError{
			Reader:  reader,
			Message: "Invalid status-code field, Expected integer value, Received: " + statusCode,
		}
	}
	SL.StatusCode = intStatusCode
	sectionsFound++

	// Intentionally lax. rfc9112.html#section-4-8
	reasonB, err := ReadUntilCRLF(reader)
	if err == io.EOF {
		return HTTPResponse{}, IncompleteInputError{
			Reader:  reader,
			Message: "Expected CRLF after status line",
		}
	}
	reason := string(reasonB)
	SL.Reason = reason
	R.StatusLine = SL
	sectionsFound++

	for !allHeadersFound {
		offsetBeforeCRLF := GetReaderOffset(reader)
		possibleHeaderLine, err := ReadUntilCRLF(reader)
		if err == io.EOF {
			if len(possibleHeaderLine) == 0 {
				return R, IncompleteInputError{
					Reader:  reader,
					Message: "Expected empty line after header section",
				}
			} else {
				return R, IncompleteInputError{
					Reader:  reader,
					Message: "Expected CRLF after header value",
				}
			}
		}
		if len(possibleHeaderLine) == 0 {
			// Headers finished
			allHeadersFound = true
			sectionsFound++
		} else {
			// We know header is present
			reader.Seek(int64(offsetBeforeCRLF), io.SeekStart)
			key, err := ReadUntil(reader, []byte(":"))
			if err == io.EOF {
				return R, IncompleteInputError{
					Reader:  reader,
					Message: "Expected ':' after header key",
				}
			}
			if key[len(key)-1] == ' ' {
				// No WS before separator
				return R, BadProtocolError{
					Reader:  reader,
					Message: "No whitespace allowed before header separator",
				}
			}
			value, err := ReadUntilCRLF(reader)
			if err == io.EOF {
				return R, IncompleteInputError{
					Reader:  reader,
					Message: "Expected CRLF after header value",
				}
			}
			H := Header{
				// 9110: 5.5-5: Replace CR, LF or NUL with SP
				Key:   string(key),
				Value: strings.TrimSpace(string(ReplaceCharsWithSpace(value, [][]byte{CR, LF, NUL}))),
			}
			R.Headers = append(R.Headers, H)
		}
	}

	// Content is present
	if R.ContentLength() != -1 {
		content, err := ReadBytes(reader, R.ContentLength())
		if err == io.EOF {
			return R, IncompleteInputError{
				Reader:  reader,
				Message: fmt.Sprintf("Expected content of length %d, Received content of length %d", R.ContentLength(), len(content)),
			}
		}
		R.Body = content
	}
	sectionsFound++

	if sectionsFound != 5 {
		return R, BadProtocolError{
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
