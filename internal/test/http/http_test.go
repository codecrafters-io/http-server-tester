package http

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var simple = []byte("GET / HTTP/1.0\r\n\r\n")

func TestParseSimple(t *testing.T) {
	hp := NewHTTPParser()

	_, err := hp.Parse(simple)
	require.NoError(t, err)

	// assert.Equal(t, n, len(simple))

	assert.Equal(t, ("HTTP/1.0"), hp.Version)

	assert.Equal(t, ("/"), hp.Path)
	assert.Equal(t, ("GET"), hp.Method)
}

func BenchmarkParseSimple(b *testing.B) {
	hp := NewHTTPParser()

	for i := 0; i < b.N; i++ {
		hp.Parse(simple)
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
	hp := NewHTTPParser()

	_, err := hp.Parse(simpleHeaders)
	require.NoError(t, err)

	assert.Equal(t, ("cookie.com"), hp.FindHeader(("Host")))
}

func BenchmarkParseSimpleHeaders(b *testing.B) {
	hp := NewHTTPParser()

	for i := 0; i < b.N; i++ {
		hp.Parse(simpleHeaders)
	}
}

var simple3Headers = []byte("GET / HTTP/1.0\r\nHost: cookie.com\r\nDate: foobar\r\nAccept: these/that\r\n\r\n")

func TestParseSimple3Headers(t *testing.T) {
	hp := NewHTTPParser()

	_, err := hp.Parse(simple3Headers)
	require.NoError(t, err)

	assert.Equal(t, ("cookie.com"), hp.FindHeader(("Host")))
	assert.Equal(t, ("foobar"), hp.FindHeader(("Date")))
	assert.Equal(t, ("these/that"), hp.FindHeader(("Accept")))
}

func BenchmarkParseSimple3Headers(b *testing.B) {
	hp := NewHTTPParser()

	for i := 0; i < b.N; i++ {
		hp.Parse(simple3Headers)
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
	hp := NewHTTPParser()

	_, err := hp.Parse(short)
	fmt.Println(err.Error())
	assert.Equal(t, err, ErrMissingData)
}

var multiline = []byte("GET / HTTP/1.0\r\nHost: cookie.com\nmore host\r\n\r\n")

func TestParseMultlineHeader(t *testing.T) {
	hp := NewHTTPParser()

	_, err := hp.Parse(multiline)
	require.NoError(t, err)

	assert.Equal(t, ("cookie.com more host"), hp.FindHeader(("Host")))
}

var specialHeaders = []byte("GET / HTTP/1.0\r\nHost: cookie.com\r\nContent-Length: 50\r\n\r\n")

func TestParseSpecialHeaders(t *testing.T) {
	hp := NewHTTPParser()

	_, err := hp.Parse(specialHeaders)
	require.NoError(t, err)

	assert.Equal(t, "cookie.com", hp.Host())
	assert.Equal(t, 50, hp.ContentLength())
}

func TestFindHeaderIgnoresCase(t *testing.T) {
	hp := NewHTTPParser()

	_, err := hp.Parse(specialHeaders)
	require.NoError(t, err)

	assert.Equal(t, ("50"), hp.FindHeader(("content-length")))
}

var multipleHeaders = []byte("GET / HTTP/1.0\r\nBar: foo\r\nBaz: quz\r\n\r\n")

func TestFindAllHeaders(t *testing.T) {
	hp := NewHTTPParser()

	_, err := hp.Parse(multipleHeaders)
	require.NoError(t, err)

	assert.Equal(t, ("foo"), hp.FindHeader("bar"))
	assert.Equal(t, ("quz"), hp.FindHeader("baz"))
}
