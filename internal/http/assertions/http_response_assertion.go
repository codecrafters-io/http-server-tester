package http_assertions

import (
	"fmt"
	"strings"

	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
)

type HTPResponseAssertion struct {
	StatusCode int    // ALWAYS REQUIRED
	Reason     string // ALWAYS REQUIRED
	Headers    []http_parser.Header
	Body       []byte
}

func NewHTTPResponseAssertion(expectedResponse http_parser.HTTPResponse) HTTPAssertion {
	return HTPResponseAssertion{StatusCode: expectedResponse.StatusLine.StatusCode, Reason: expectedResponse.StatusLine.Reason, Headers: expectedResponse.Headers, Body: expectedResponse.Body}
}

func (a HTPResponseAssertion) Run(response http_parser.HTTPResponse) error {
	actualStatusLine := response.StatusLine
	if actualStatusLine.StatusCode != a.StatusCode {
		return fmt.Errorf("Expected status code %d, got %d", a.StatusCode, actualStatusLine.StatusCode)
	}

	if actualStatusLine.Reason != a.Reason {
		return fmt.Errorf("Expected reason %s, got %s", a.Reason, actualStatusLine.Reason)
	}

	if response.Headers != nil {
		if len(response.Headers) != len(a.Headers) {
			return fmt.Errorf("Expected %d headers, got %d", len(a.Headers), len(response.Headers))
		}
		for _, header := range a.Headers {
			expectedKey, expectedValue := header.Key, header.Value
			actualValue := response.FindHeader(expectedKey)
			if !strings.EqualFold(actualValue, expectedValue) {
				return fmt.Errorf("Expected header %s: %s, got %s", expectedKey, expectedValue, actualValue)
			}
		}
	}

	if response.Body != nil {
		if len(response.Body) != len(a.Body) {
			return fmt.Errorf("Expected body of length %d, got %d", len(a.Body), len(response.Body))
		}
		if string(response.Body) != string(a.Body) {
			return fmt.Errorf("Expected body %s, got %s", a.Body, response.Body)
		}
	}
	return nil
}
