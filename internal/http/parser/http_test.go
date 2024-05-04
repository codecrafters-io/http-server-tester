package http_parser

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var simpleResponse = []byte("HTTP/1.0 200 OK\r\n\r\n")

func TestParseSimpleResponse(t *testing.T) {
	response, n, err := Parse(simpleResponse)

	require.NoError(t, err)
	assert.Equal(t, n, len(simpleResponse))

	assert.Equal(t, "HTTP/1.0", response.StatusLine.Version)

	assert.Equal(t, 200, response.StatusLine.StatusCode)
	assert.Equal(t, "OK", response.StatusLine.Reason)
}

var simpleHeaders = []byte("HTTP/1.0 200 OK\r\nServer: Werkzeug/3.0.2 Python/3.10.13\r\n\r\n")

func TestParseSimpleHeaders(t *testing.T) {
	response, _, err := Parse(simpleHeaders)
	require.NoError(t, err)

	assert.Equal(t, "Werkzeug/3.0.2 Python/3.10.13", response.FindHeader("Server"))
}

var multipleHeaders = []byte("HTTP/1.0 200 OK\r\nServer: Werkzeug/3.0.2 Python/3.10.13\r\nDate: Tue, 30 Apr 2024 06:16:31 GMT\r\n\r\n")

func TestParseMultiHeaders(t *testing.T) {
	response, _, err := Parse(multipleHeaders)
	require.NoError(t, err)

	assert.Equal(t, "Werkzeug/3.0.2 Python/3.10.13", response.FindHeader("Server"))
	assert.Equal(t, "Tue, 30 Apr 2024 06:16:31 GMT", response.FindHeader("Date"))
}

var specialHeaders = []byte("HTTP/1.0 200 OK\r\nServer: Werkzeug/3.0.2 Python/3.10.13\r\nDate: Tue, 30 Apr 2024 06:16:31 GMT\r\nContent-Length: 0\r\n\r\n")

func TestParseSpecialHeaders(t *testing.T) {
	response, _, err := Parse(specialHeaders)
	require.NoError(t, err)

	assert.Equal(t, "0", response.FindHeader("Content-Length"))
	assert.Equal(t, "0", response.FindHeader("content-length"))
	assert.Equal(t, 0, response.ContentLength())
}

var multiline = []byte("HTTP/1.0 200 OK\r\nHost: cookie.com\nmore host\r\n\r\n")

func TestParseMultilineHeader(t *testing.T) {
	response, _, err := Parse(multiline)
	require.NoError(t, err)

	assert.Equal(t, "cookie.com more host", response.FindHeader("Host"))
}

func TestFindHeaderIgnoresCase(t *testing.T) {
	response, _, err := Parse(specialHeaders)
	require.NoError(t, err)

	assert.Equal(t, "0", response.FindHeader("content-length"))
}

var missingContent = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain;charset=utf-8\r\nContent-Length: 29\r\nConnection: close\r\n\r\nThis is a plain text ")

func TestParseResponseMissingContent(t *testing.T) {
	_, _, err := Parse(missingContent)

	expected := strings.TrimSpace(`
Received: "his is a plain text "
                               ^ error
Error: Expected content of length 29, Received content of length 21`)

	assert.Error(t, err)
	assert.IsType(t, IncompleteInputError{}, err)
	assert.EqualError(t, err, expected)
}

var properContent = []byte("HTTP/1.1 200 OK\r\nServer: Werkzeug/3.0.2 Python/3.10.13\r\nDate: Tue, 30 Apr 2024 06:16:31 GMT\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Length: 29\r\nConnection: close\r\n\r\nThis is a plain text response")

func TestParseResponseProperContent(t *testing.T) {
	response, n, err := Parse(properContent)
	require.NoError(t, err)
	assert.Equal(t, n, len(properContent))

	assert.Equal(t, "HTTP/1.1", response.StatusLine.Version)

	assert.Equal(t, 200, response.StatusLine.StatusCode)
	assert.Equal(t, "OK", response.StatusLine.Reason)

	assert.Equal(t, "Werkzeug/3.0.2 Python/3.10.13", response.FindHeader("Server"))
	assert.Equal(t, "Tue, 30 Apr 2024 06:16:31 GMT", response.FindHeader("Date"))
	assert.Equal(t, "29", response.FindHeader("Content-Length"))

	assert.Equal(t, "This is a plain text response", string(response.Body))
}

var protocolErrors = map[string]string{
	"HTTP/3.0 200 OK": `
Received: "HTTP/3.0 200 OK"
                ^ error
Error: Expected HTTP-version 1.X, Received: 3.0`,
	"HTPP/1.0 200 OK": `
Received: "HTPP/1.0 200 OK"
               ^ error
Error: Expected HTTP-version field to start with 'HTTP/'`,
	"HTTP|1.0 200 OK": `
Received: "HTTP|1.0 200 OK"
               ^ error
Error: Expected HTTP-version field to start with 'HTTP/'`,
	"HTTP//1.0 200 OK": `
Received: "HTTP//1.0 200 OK"
                   ^ error
Error: Invalid HTTP-version field length`,
	"HTTP/1 200 OK": `
Received: "HTTP/1 200 OK"
                ^ error
Error: Invalid HTTP-version field length`,
	"HTTP/1.0 ZOO OK": `
Received: "HTTP/1.0 ZOO OK"
                      ^ error
Error: Invalid status-code field, Expected integer value, Received: ZOO`,
	"HTTP/1.0 2000 OK": `
Received: "HTTP/1.0 2000 OK"
                       ^ error
Error: Invalid status-code field length, Expected: 3 digits, Received: 4`,
	"HTTP/1.1 200 OK\r\nConnection : close\r\n": `
Received: "200 OK\r\nConnection : close\r\n"
                                 ^ error
Error: No whitespace allowed before header separator`,
}

func TestParseVersion(t *testing.T) {
	for response, error := range protocolErrors {
		t.Run(response, func(t *testing.T) {
			_, _, err := Parse([]byte(response))
			assert.Error(t, err)
			assert.IsType(t, BadProtocolError{}, err)
			assert.EqualError(t, err, strings.TrimSpace(error))
		})
	}
}

var incompleteErrors = map[string]string{
	"HTTP/1.0_200_OK": `
Received: "HTTP/1.0_200_OK"
                          ^ error
Error: Expected: HTTP-version, Received: "HTTP/1.0_200_OK"`,
	"HTTP/1.0 200 OK": `
Received: "HTTP/1.0 200 OK"
                          ^ error
Error: Expected CRLF after status line`,
	"HTTP/1.1 200 OK\r\nConnection: close": `
Received: "K\r\nConnection: close"
                                 ^ error
Error: Expected CRLF after header value`,
	"HTTP/1.1 200 OK\r\nConnection: close\r\n": `
Received: "\nConnection: close\r\n"
                                  ^ error
Error: Expected empty line after header section`,
	"HTTP/1.1 200 OK\r\nConnection\r\n": `
Received: "200 OK\r\nConnection\r\n"
                                   ^ error
Error: Expected ':' after header key`,
	"HTTP/1.1 200\r\n": `
Received: "HTTP/1.1 200\r\n"
                           ^ error
Error: Status line has missing sections, Expected: HTTP-version status-code reason-phrase`,
}

func TestParseVersion2(t *testing.T) {
	for response, error := range incompleteErrors {
		t.Run(response, func(t *testing.T) {
			_, _, err := Parse([]byte(response))
			assert.Error(t, err)
			assert.IsType(t, IncompleteInputError{}, err)
			assert.EqualError(t, err, strings.TrimSpace(error))
		})
	}
}

func BenchmarkParseSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := Parse(simpleResponse)
		if err != nil {
			return
		}
	}
}

func BenchmarkNetHTTP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bufio.NewReader(bytes.NewReader(simpleResponse))
		_, err := http.ReadRequest(buf)
		if err != nil {
			return
		}
	}
}
func BenchmarkParseSimpleHeaders(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := Parse(simpleHeaders)
		if err != nil {
			return
		}
	}
}

func BenchmarkParseMultiHeaders(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := Parse(multipleHeaders)
		if err != nil {
			return
		}
	}
}

func BenchmarkNetHTTP3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bufio.NewReader(bytes.NewReader(multipleHeaders))
		_, err := http.ReadRequest(buf)
		if err != nil {
			return
		}
	}
}
