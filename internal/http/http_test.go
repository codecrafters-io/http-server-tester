package http

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"
	"testing"

	http_request "github.com/codecrafters-io/http-server-tester/internal/http/parser/request"
	http_response "github.com/codecrafters-io/http-server-tester/internal/http/parser/response"
	http_utils "github.com/codecrafters-io/http-server-tester/internal/http/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var simpleRequest = []byte("GET / HTTP/1.0\r\n\r\n")

func TestParseSimpleRequest(t *testing.T) {
	request, n, err := http_request.Parse(simpleRequest)

	require.NoError(t, err)
	assert.Equal(t, n, len(simpleRequest))

	assert.Equal(t, ("HTTP/1.0"), request.RequestLine.Version)

	assert.Equal(t, ("/"), request.RequestLine.Path)
	assert.Equal(t, ("GET"), request.RequestLine.Method)
}

var simpleHeaders = []byte("GET / HTTP/1.0\r\nHost: cookie.com\r\n\r\n")

func TestParseSimpleHeaders(t *testing.T) {
	request, _, err := http_request.Parse(simpleHeaders)
	require.NoError(t, err)

	assert.Equal(t, ("cookie.com"), request.FindHeader(("Host")))
}

var multipleHeaders = []byte("GET / HTTP/1.0\r\nHost: cookie.com\r\nDate: foobar\r\nAccept: these/that\r\n\r\n")

func TestParseMultiHeaders(t *testing.T) {
	request, _, err := http_request.Parse(multipleHeaders)
	require.NoError(t, err)

	assert.Equal(t, ("cookie.com"), request.FindHeader(("Host")))
	assert.Equal(t, ("foobar"), request.FindHeader(("Date")))
	assert.Equal(t, ("these/that"), request.FindHeader(("Accept")))
}

var short = []byte("GET / HT")

func TestParseMissingData(t *testing.T) {
	_, _, err := http_request.Parse(short)

	expected := strings.TrimSpace(`
Received: "GET / HT"
                   ^ error
Error: Expected CRLF after status line`)

	assert.Error(t, err)
	assert.IsType(t, http_utils.IncompleteInputError{}, err)
	assert.EqualError(t, err, expected)
}

var multiline = []byte("GET / HTTP/1.0\r\nHost: cookie.com\nmore host\r\n\r\n")

func TestParseMultlineHeader(t *testing.T) {
	request, _, err := http_request.Parse(multiline)
	require.NoError(t, err)

	assert.Equal(t, ("cookie.com more host"), request.FindHeader(("Host")))
}

var specialHeaders = []byte("GET / HTTP/1.0\r\nHost: cookie.com\r\nContent-Length: 0\r\n\r\n")

func TestParseSpecialHeaders(t *testing.T) {
	request, _, err := http_request.Parse(specialHeaders)
	require.NoError(t, err)

	assert.Equal(t, "cookie.com", request.Host())
	assert.Equal(t, 0, request.ContentLength())
}

func TestFindHeaderIgnoresCase(t *testing.T) {
	request, _, err := http_request.Parse(specialHeaders)
	require.NoError(t, err)

	assert.Equal(t, ("0"), request.FindHeader(("content-length")))
}

var simpleResponse = []byte("HTTP/1.0 200 OK\r\n\r\n")

func TestParseSimpleResponse(t *testing.T) {
	response, n, err := http_response.Parse(simpleResponse)

	require.NoError(t, err)
	assert.Equal(t, n, len(simpleResponse))

	assert.Equal(t, ("HTTP/1.0"), response.StatusLine.Version)

	assert.Equal(t, (200), response.StatusLine.StatusCode)
	assert.Equal(t, ("OK"), response.StatusLine.Reason)
}

var missingData = []byte("HTTP/1.1 200\r\n")

func TestParseResponseMissingData(t *testing.T) {
	_, _, err := http_response.Parse(missingData)

	expected := strings.TrimSpace(`
Received: "HTTP/1.1 200\r\n"
                           ^ error
Error: Expected SP between elements of the status line`)

	assert.Error(t, err)
	assert.IsType(t, http_utils.IncompleteInputError{}, err)
	assert.EqualError(t, err, expected)
}

var missingContent = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain;charset=utf-8\r\nContent-Length: 29\r\nConnection: close\r\n\r\nThis is a plain text ")

var properContent = []byte("HTTP/1.1 200 OK\r\nServer: Werkzeug/3.0.2 Python/3.10.13\r\nDate: Tue, 30 Apr 2024 06:16:31 GMT\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Length: 29\r\nConnection: close\r\n\r\nThis is a plain text response")

func TestParseResponseMissingContent(t *testing.T) {
	_, _, err := http_response.Parse(missingContent)

	expected := strings.TrimSpace(`
Received: "his is a plain text "
                               ^ error
Error: Expected content of length 29`)

	assert.Error(t, err)
	assert.IsType(t, http_utils.IncompleteInputError{}, err)
	assert.EqualError(t, err, expected)
}

func TestParseResponseProperContent(t *testing.T) {
	response, n, err := http_response.Parse(properContent)
	require.NoError(t, err)
	assert.Equal(t, n, len(properContent))

	assert.Equal(t, ("HTTP/1.1"), response.StatusLine.Version)

	assert.Equal(t, (200), response.StatusLine.StatusCode)
	assert.Equal(t, ("OK"), response.StatusLine.Reason)

	assert.Equal(t, ("Werkzeug/3.0.2 Python/3.10.13"), response.FindHeader(("Server")))
	assert.Equal(t, ("Tue, 30 Apr 2024 06:16:31 GMT"), response.FindHeader(("Date")))
	assert.Equal(t, ("29"), response.FindHeader(("Content-Length")))
}

func BenchmarkParseSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		http_request.Parse(simpleRequest)
	}
}

func BenchmarkNetHTTP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bufio.NewReader(bytes.NewReader(simpleRequest))
		http.ReadRequest(buf)
	}
}
func BenchmarkParseSimpleHeaders(b *testing.B) {
	for i := 0; i < b.N; i++ {
		http_request.Parse(simpleHeaders)
	}
}

func BenchmarkParseMultiHeaders(b *testing.B) {
	for i := 0; i < b.N; i++ {
		http_request.Parse(multipleHeaders)
	}
}

func BenchmarkNetHTTP3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bufio.NewReader(bytes.NewReader(multipleHeaders))
		http.ReadRequest(buf)
	}
}
