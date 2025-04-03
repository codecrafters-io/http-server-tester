package internal

import (
	"fmt"
	"net/http"

	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
)

type RequestResponsePair struct {
	Request  *http.Request
	Response *http_parser.HTTPResponse
}

// Base: GET /

func getBaseURLGetRequest() (*http.Request, error) {
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create request: %v", err)
	}
	return request, nil
}

func getBaseURLGetResponse() (*http_parser.HTTPResponse, error) {
	response := http_parser.HTTPResponse{
		StatusLine: http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"},
	}
	return &response, nil
}

func GetBaseURLGetRequestResponsePair() (*RequestResponsePair, error) {
	request, err := getBaseURLGetRequest()
	if err != nil {
		return nil, err
	}
	response, err := getBaseURLGetResponse()
	if err != nil {
		return nil, err
	}
	return &RequestResponsePair{Request: request, Response: response}, nil
}

// Echo: GET /echo/{content}

func getEchoRequest(content string) (*http.Request, error) {
	request, err := http.NewRequest("GET", URL+"echo/"+content, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create request: %v", err)
	}
	return request, nil
}

func getEchoResponse(content string) (*http_parser.HTTPResponse, error) {
	statusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}

	header1 := http_parser.Header{Key: "Content-Type", Value: "text/plain"}
	header2 := http_parser.Header{Key: "Content-Length", Value: fmt.Sprintf("%d", len(content))}
	headers := []http_parser.Header{header1, header2}

	body := []byte(content)

	response := http_parser.HTTPResponse{StatusLine: statusLine, Headers: headers, Body: body}

	return &response, nil
}

func getEchoRequestResponsePair(content string) (*RequestResponsePair, error) {
	request, err := getEchoRequest(content)
	if err != nil {
		return nil, err
	}
	response, err := getEchoResponse(content)
	if err != nil {
		return nil, err
	}
	return &RequestResponsePair{Request: request, Response: response}, nil
}

func GetEchoRequestResponsePair() (*RequestResponsePair, error) {
	return getEchoRequestResponsePair(randomUrlPath())
}

// User-Agent: GET /user-agent

func getUserAgentRequest(userAgent string) (*http.Request, error) {
	request, err := http.NewRequest("GET", URL+"user-agent", nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create request: %v", err)
	}
	request.Header.Set("User-Agent", userAgent)

	return request, nil
}

func getUserAgentResponse(userAgent string) (*http_parser.HTTPResponse, error) {
	statusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}

	header1 := http_parser.Header{Key: "Content-Type", Value: "text/plain"}
	header2 := http_parser.Header{Key: "Content-Length", Value: fmt.Sprintf("%d", len(userAgent))}
	headers := []http_parser.Header{header1, header2}
	body := []byte(userAgent)

	response := http_parser.HTTPResponse{StatusLine: statusLine, Headers: headers, Body: body}
	return &response, nil
}

func getUserAgentRequestResponsePair(userAgent string) (*RequestResponsePair, error) {
	request, err := getUserAgentRequest(userAgent)
	if err != nil {
		return nil, err
	}
	response, err := getUserAgentResponse(userAgent)
	if err != nil {
		return nil, err
	}
	return &RequestResponsePair{Request: request, Response: response}, nil
}

func GetUserAgentRequestResponsePair() (*RequestResponsePair, error) {
	return getUserAgentRequestResponsePair(randomUserAgent())
}
