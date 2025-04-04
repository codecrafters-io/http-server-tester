package internal

import (
	"fmt"
	"net/http"

	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
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

// files: GET /files/{filename}

func getFilesRequest(filename, fileContent string, logger *logger.Logger) (*http.Request, error) {
	// Prerequisite: create file
	logger.Infof("Testing existing file")
	logger.Debugf("Creating file %s in %s", filename, DATA_DIR)
	logger.Debugf("File Content: %q", fileContent)
	err := createFileWith(DATA_DIR+filename, fileContent)
	if err != nil {
		return nil, fmt.Errorf("Could not create file: %v", err)
	}

	request, err := http.NewRequest("GET", URL+"files/"+filename, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create request: %v", err)
	}

	return request, nil
}

func getFilesResponse(fileContent string) (*http_parser.HTTPResponse, error) {
	statusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}

	header1 := http_parser.Header{Key: "Content-Type", Value: "application/octet-stream"}
	header2 := http_parser.Header{Key: "Content-Length", Value: fmt.Sprintf("%d", len(fileContent))}
	headers := []http_parser.Header{header1, header2}
	body := []byte(fileContent)

	response := http_parser.HTTPResponse{StatusLine: statusLine, Headers: headers, Body: body}
	return &response, nil
}

func getFilesRequestResponsePair(filename, fileContent string, logger *logger.Logger) (*RequestResponsePair, error) {
	request, err := getFilesRequest(filename, fileContent, logger)
	if err != nil {
		return nil, err
	}
	response, err := getFilesResponse(fileContent)
	if err != nil {
		return nil, err
	}
	return &RequestResponsePair{Request: request, Response: response}, nil
}

func GetFilesRequestResponsePair(logger *logger.Logger) (*RequestResponsePair, error) {
	return getFilesRequestResponsePair(randomFileName(), randomFileContent(), logger)
}

func getRandomRequestResponsePair(logger *logger.Logger) (*RequestResponsePair, error) {
	countOfPossibleRequestResponsePairs := 4

	possibleRequestResponsePairs := []func() (*RequestResponsePair, error){
		GetBaseURLGetRequestResponsePair,
		GetEchoRequestResponsePair,
		GetUserAgentRequestResponsePair,
		// GetFilesRequestResponsePair, // Expected mismatch in interface
	}

	randomIndex := random.RandomInt(0, countOfPossibleRequestResponsePairs)

	if randomIndex == 3 {
		requestResponsePair, err := GetFilesRequestResponsePair(logger)
		if err != nil {
			return nil, err
		}
		return requestResponsePair, nil
	}

	requestResponsePair, err := possibleRequestResponsePairs[randomIndex]()
	if err != nil {
		return nil, err
	}
	return requestResponsePair, nil
}

// GetRandomRequestResponsePairs returns a slice of RequestResponsePairs
// The RequestResponsePairs can be of the following types:
// - GET /
// - GET /echo/{content}
// - GET /user-agent
// - GET /files/{filename}
// Use with Data Directory, and pass --directory flag to the server
func GetRandomRequestResponsePairs(count int, logger *logger.Logger) ([]*RequestResponsePair, error) {
	requestResponsePairs := make([]*RequestResponsePair, count)
	for i := range count {
		requestResponsePair, err := getRandomRequestResponsePair(logger)
		if err != nil {
			return nil, err
		}
		requestResponsePairs[i] = requestResponsePair
	}
	return requestResponsePairs, nil
}
