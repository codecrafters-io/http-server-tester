package internal

import (
	"net/http"
	"strings"
	"testing"

	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
)

func TestCurlCommandGeneration(t *testing.T) {
	testCases := []struct {
		name       string
		method     string
		url        string
		headers    map[string]string
		body       string
		curlOutput string
	}{
		{
			name:       "GET request with headers",
			method:     "GET",
			url:        "https://example.com",
			headers:    map[string]string{"Header1": "Value1", "Header2": "Value2"},
			body:       "",
			curlOutput: "curl -v -X GET https://example.com -H \"Header1: Value1\" -H \"Header2: Value2\"",
		},
		{
			name:       "POST request with body and headers",
			method:     "POST",
			url:        "https://example.com",
			headers:    map[string]string{"Content-Type": "application/json", "Authorization": "Bearer Token"},
			body:       `{"key": "value"}`,
			curlOutput: "curl -v -X POST https://example.com -H \"Authorization: Bearer Token\" -H \"Content-Type: application/json\" -d '{\"key\": \"value\"}'",
		},
		{
			name:       "Empty body",
			method:     "PUT",
			url:        "https://example.com/resource",
			headers:    map[string]string{"Content-Type": "text/plain"},
			body:       "",
			curlOutput: "curl -v -X PUT https://example.com/resource -H \"Content-Type: text/plain\"",
		},
		{
			name:       "No headers",
			method:     "DELETE",
			url:        "https://example.com/item",
			headers:    map[string]string{},
			body:       "request body",
			curlOutput: "curl -v -X DELETE https://example.com/item -d 'request body'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			curlCommand := http_connection.HttpRequestToCurlString(req)

			if curlCommand != tc.curlOutput {
				t.Errorf("Expected: %s\nActual: %s", tc.curlOutput, curlCommand)
			}
		})
	}
}

func stringToReadCloser(s string) *StringReadCloser {
	return &StringReadCloser{strings.NewReader(s)}
}

type StringReadCloser struct {
	*strings.Reader
}

func (src *StringReadCloser) Close() error { return nil }
