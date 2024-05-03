package http

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"testing"

	http_request "github.com/codecrafters-io/http-server-tester/internal/http/parser/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var simple = []byte("GET / HTTP/1.0\r\n\r\n")

func TestParseSimple(t *testing.T) {
	request, n, err := http_request.Parse(simple)

	require.NoError(t, err)
	assert.Equal(t, n, len(simple))

	assert.Equal(t, ("HTTP/1.0"), request.RequestLine.Version)

	assert.Equal(t, ("/"), request.RequestLine.Path)
	assert.Equal(t, ("GET"), request.RequestLine.Method)
}

func BenchmarkParseSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		http_request.Parse(simple)
	}
}

func BenchmarkNetHTTP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bufio.NewReader(bytes.NewReader(simple))
		http.ReadRequest(buf)
	}
}

var simpleHeaders = []byte("GET / HTTP/1.0\r\nHost: cookie.com\r\n\r\n")

func TestParseSimpleHeaders(t *testing.T) {
	request, _, err := http_request.Parse(simpleHeaders)
	require.NoError(t, err)

	assert.Equal(t, ("cookie.com"), request.FindHeader(("Host")))
}

func BenchmarkParseSimpleHeaders(b *testing.B) {
	for i := 0; i < b.N; i++ {
		http_request.Parse(simpleHeaders)
	}
}

var simple3Headers = []byte("GET / HTTP/1.0\r\nHost: cookie.com\r\nDate: foobar\r\nAccept: these/that\r\n\r\n")

func TestParseSimple3Headers(t *testing.T) {
	request, _, err := http_request.Parse(simple3Headers)
	require.NoError(t, err)

	assert.Equal(t, ("cookie.com"), request.FindHeader(("Host")))
	assert.Equal(t, ("foobar"), request.FindHeader(("Date")))
	assert.Equal(t, ("these/that"), request.FindHeader(("Accept")))
}

func BenchmarkParseSimple3Headers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		http_request.Parse(simple3Headers)
	}
}

func BenchmarkNetHTTP3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bufio.NewReader(bytes.NewReader(simple3Headers))
		http.ReadRequest(buf)
	}
}

var short = []byte("GET / HT")

func TestParseMissingData(t *testing.T) {
	_, _, err := http_request.Parse(short)
	fmt.Println(err.Error())
	assert.Equal(t, err, http_request.ErrMissingData)
}

var multiline = []byte("GET / HTTP/1.0\r\nHost: cookie.com\nmore host\r\n\r\n")

func TestParseMultlineHeader(t *testing.T) {
	request, _, err := http_request.Parse(multiline)
	require.NoError(t, err)

	assert.Equal(t, ("cookie.com more host"), request.FindHeader(("Host")))
}

var specialHeaders = []byte("GET / HTTP/1.0\r\nHost: cookie.com\r\nContent-Length: 50\r\n\r\n")

func TestParseSpecialHeaders(t *testing.T) {
	request, _, err := http_request.Parse(specialHeaders)
	require.NoError(t, err)

	assert.Equal(t, "cookie.com", request.Host())
	assert.Equal(t, 50, request.ContentLength())
}

func TestFindHeaderIgnoresCase(t *testing.T) {
	request, _, err := http_request.Parse(specialHeaders)
	require.NoError(t, err)

	assert.Equal(t, ("50"), request.FindHeader(("content-length")))
}

var multipleHeaders = []byte("GET / HTTP/1.0\r\nBar: foo\r\nBaz: quz\r\n\r\n")

func TestFindAllHeaders(t *testing.T) {
	request, _, err := http_request.Parse(multipleHeaders)
	require.NoError(t, err)

	assert.Equal(t, ("foo"), request.FindHeader("bar"))
	assert.Equal(t, ("quz"), request.FindHeader("baz"))
}
