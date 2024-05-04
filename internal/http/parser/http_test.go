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

var short = []byte("HTTP/1 200 OK")

func TestParseMissingData(t *testing.T) {
	_, _, err := Parse(short)

	expected := strings.TrimSpace(`
Received: "HTTP/1 200 OK"
                ^ error
Error: Invalid HTTP-version field length`)

	assert.Error(t, err)
	assert.IsType(t, BadProtocolError{}, err)
	assert.EqualError(t, err, expected)
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

var missingData = []byte("HTTP/1.1 200\r\n")

func TestParseResponseMissingData(t *testing.T) {
	_, _, err := Parse(missingData)

	expected := strings.TrimSpace(`
Received: "HTTP/1.1 200\r\n"
                           ^ error
Error: Status line has missing sections, Expected: HTTP-version status-code reason-phrase`)

	assert.Error(t, err)
	assert.IsType(t, IncompleteInputError{}, err)
	assert.EqualError(t, err, expected)
}

var missingContent = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain;charset=utf-8\r\nContent-Length: 29\r\nConnection: close\r\n\r\nThis is a plain text ")

var properContent = []byte("HTTP/1.1 200 OK\r\nServer: Werkzeug/3.0.2 Python/3.10.13\r\nDate: Tue, 30 Apr 2024 06:16:31 GMT\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Length: 29\r\nConnection: close\r\n\r\nThis is a plain text response")

func TestParseResponseMissingContent(t *testing.T) {
	_, _, err := Parse(missingContent)

	expected := strings.TrimSpace(`
Received: "his is a plain text "
                               ^ error
Error: Expected content of length 29`)
	// ToDo: Expected: X received: Y

	assert.Error(t, err)
	assert.IsType(t, IncompleteInputError{}, err)
	assert.EqualError(t, err, expected)
}

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

var wrongHTTP0 = []byte("HTTP/3.0 200 OK")
var wrongHTTP1 = []byte("HTPP/1.0 200 OK")
var wrongHTTP2 = []byte("HTTP|1.0 200 OK")
var wrongHTTP3 = []byte("HTTP//1.0 200 OK")
var wrongHTTP4 = []byte("HTTP/1 200 OK")
var wrongHTTP5 = []byte("HTTP/1.0 ZOO OK")
var wrongHTTP6 = []byte("HTTP/1.0 2000 OK")
var wrongHTTP7 = []byte("HTTP/1.0_200_OK")
var wrongHTTP8 = []byte("HTTP/1.0 200 OK")

var wrongHTTP = map[string]string{`
Received: "HTTP/3.0 200 OK"
                ^ error
Error: Expected HTTP-version 1.X, Received: 3.0`: "HTTP/3.0 200 OK", `
Received: "HTPP/1.0 200 OK"
               ^ error
Error: Expected HTTP-version field to start with 'HTTP/'`: "HTPP/1.0 200 OK",
}

var wrong = map[string]string{
	"HTTP/3.0 200 OK": `
Received: "HTTP/3.0 200 OK"
                ^ error
Error: Expected HTTP-version 1.X, Received: 3.0`,
	"HTPP/1.0 200 OK": `
Received: "HTPP/1.0 200 OK"
               ^ error
Error: Expected HTTP-version field to start with 'HTTP/'`,
	"HTTP|1.0 200 OK": `
    `,
	"HTTP//1.0 200 OK": `
    `,
	"HTTP/1 200 OK": `
    `,
	"HTTP/1.0 ZOO OK": `
    `,
	"HTTP/1.0 2000 OK": `
    `,
	"HTTP/1.0_200_OK": `
    `,
	"HTTP/1.0 200 OK": `
    `,
}

func TestParseVersion(t *testing.T) {
	for response, error := range wrong {
		t.Run(response, func(t *testing.T) {
			_, _, err := Parse([]byte(response))
			assert.Error(t, err)
			assert.IsType(t, BadProtocolError{}, err)
			assert.EqualError(t, err, strings.TrimSpace(error))
		})
	}
}

func TestParseVersion1(t *testing.T) {
	_, _, err := Parse(wrongHTTP0)

	expected := strings.TrimSpace(`
Received: "HTTP/3.0 200 OK"
                ^ error
Error: Expected HTTP-version 1.X, Received: 3.0`)

	assert.Error(t, err)
	assert.IsType(t, BadProtocolError{}, err)
	assert.EqualError(t, err, expected)
}

func TestParseVersion2(t *testing.T) {
	_, _, err := Parse(wrongHTTP1)
	expected := strings.TrimSpace(`
Received: "HTPP/1.0 200 OK"
               ^ error
Error: Expected HTTP-version field to start with 'HTTP/'`)
	assert.Error(t, err)
	assert.IsType(t, BadProtocolError{}, err)
	assert.EqualError(t, err, expected)
}

func TestParseVersion3(t *testing.T) {
	_, _, err := Parse(wrongHTTP2)
	expected := strings.TrimSpace(`
Received: "HTTP|1.0 200 OK"
               ^ error
Error: Expected HTTP-version field to start with 'HTTP/'`)
	assert.Error(t, err)
	assert.IsType(t, BadProtocolError{}, err)
	assert.EqualError(t, err, expected)
}

func TestParseVersion4(t *testing.T) {
	_, _, err := Parse(wrongHTTP3)
	expected := strings.TrimSpace(`
Received: "HTTP//1.0 200 OK"
                   ^ error
Error: Invalid HTTP-version field length`)
	assert.Error(t, err)
	assert.IsType(t, BadProtocolError{}, err)
	assert.EqualError(t, err, expected)
}

func TestParseVersion5(t *testing.T) {
	_, _, err := Parse(wrongHTTP4)
	expected := strings.TrimSpace(`
Received: "HTTP/1 200 OK"
                ^ error
Error: Invalid HTTP-version field length`)
	assert.Error(t, err)
	assert.IsType(t, BadProtocolError{}, err)
	assert.EqualError(t, err, expected)
}

func TestParseVersion6(t *testing.T) {
	_, _, err := Parse(wrongHTTP5)
	expected := strings.TrimSpace(`
Received: "HTTP/1.0 ZOO OK"
                      ^ error
Error: Invalid status-code field, Expected integer value, Received: ZOO`)
	assert.Error(t, err)
	assert.IsType(t, BadProtocolError{}, err)
	assert.EqualError(t, err, expected)
}

func TestParseVersion7(t *testing.T) {
	_, _, err := Parse(wrongHTTP6)
	expected := strings.TrimSpace(`
Received: "HTTP/1.0 2000 OK"
                       ^ error
Error: Invalid status-code field length, Expected: 3 digits, Received: 4`)
	assert.Error(t, err)
	assert.IsType(t, BadProtocolError{}, err)
	assert.EqualError(t, err, expected)
}

func TestParseVersion8(t *testing.T) {
	_, _, err := Parse(wrongHTTP7)
	expected := strings.TrimSpace(`
Received: "HTTP/1.0_200_OK"
                          ^ error
Error: Expected: HTTP-version, Received: "HTTP/1.0_200_OK"`)
	assert.Error(t, err)
	assert.IsType(t, IncompleteInputError{}, err)
	assert.EqualError(t, err, expected)
}

func TestParseVersion9(t *testing.T) {
	_, _, err := Parse(wrongHTTP8)
	expected := strings.TrimSpace(`
Received: "HTTP/1.0 200 OK"
                          ^ error
Error: Expected CRLF after status line`)
	assert.Error(t, err)
	assert.IsType(t, IncompleteInputError{}, err)
	assert.EqualError(t, err, expected)
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
