package http_assertions

import (
	"fmt"
	"strings"

	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
)

type HTTPResponseAssertion struct {
	StatusCode int    // ALWAYS REQUIRED
	Reason     string // ALWAYS REQUIRED
	Headers    []http_parser.Header
	Body       []byte
}

func NewHTTPResponseAssertion(expectedResponse http_parser.HTTPResponse) HTTPResponseAssertion {
	return HTTPResponseAssertion{StatusCode: expectedResponse.StatusLine.StatusCode, Reason: expectedResponse.StatusLine.Reason, Headers: expectedResponse.Headers, Body: expectedResponse.Body}
}

func (a HTTPResponseAssertion) Run(response http_parser.HTTPResponse) error {
	actualStatusLine := response.StatusLine
	if actualStatusLine.StatusCode != a.StatusCode {
		return fmt.Errorf("Expected status code %d, got %d", a.StatusCode, actualStatusLine.StatusCode)
	}

	if actualStatusLine.Reason != a.Reason {
		return fmt.Errorf("Expected reason to be %q, got %q", a.Reason, actualStatusLine.Reason)
	}

	if a.Headers != nil {
		// Only if we pass Headers in the HTTPResponseAssertion, we will check the headers
		for _, header := range a.Headers {
			expectedKey, expectedValue := header.Key, header.Value
			actualValue := response.FindHeader(expectedKey)
			if !strings.EqualFold(actualValue, expectedValue) {
				return fmt.Errorf("Expected header %s: %s, got %s", expectedKey, expectedValue, actualValue)
			}
		}
	}

	if a.Body != nil {
		// Only if we pass Body in the HTTPResponseAssertion, we will check the body
		if len(response.Body) != len(a.Body) {
			return fmt.Errorf("Expected body of length %d, got %d", len(a.Body), len(response.Body))
		}
		if string(response.Body) != string(a.Body) {
			return fmt.Errorf("Expected body %s, got %s", a.Body, response.Body)
		}
	}
	return nil
}
