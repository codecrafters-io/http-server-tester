package internal

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"

	logger "github.com/codecrafters-io/tester-utils/logger"
)

const URL = "http://localhost:4221/"
const TCP_DEST = "localhost:4221"
const DATA_DIR = "/tmp/data/codecrafters.io/http-server-tester/"
const FILENAME_SIZE = 40

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getFirstLine(s string) string {
	lines := strings.Split(s, "\r\n")
	if len(lines) == 0 {
		return ""
	}
	return lines[0]
}

func logFriendlyHTTPMessage(logger *logger.Logger, msg string, logPrefix string) {
	for _, line := range strings.Split(msg, "\r\n") {
		logger.Debugf("%s %s", logPrefix, line)
	}
}

func dumpRequest(logger *logger.Logger, req *http.Request) error {
	logCurl(logger, req)
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return fmt.Errorf("Failed to dump request: '%v'", err)
	}
	logger.Infof("Sending request (status line): %s", getFirstLine(string(reqDump)))
	logPrefix := ">>>"
	logger.Debugf("Sending request: (Messages with %s prefix are part of this log)", logPrefix)
	logFriendlyHTTPMessage(logger, string(reqDump), logPrefix)

	return nil
}

func dumpResponse(logger *logger.Logger, resp *http.Response) error {
	respDump, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return fmt.Errorf("Failed to dump response: '%v'", err)
	}
	logPrefix := ">>>"
	logger.Debugf("Received response: (Messages with %s prefix are part of this log)", logPrefix)
	logFriendlyHTTPMessage(logger, string(respDump), logPrefix)

	return nil
}

func dumpResponseWithBody(logger *logger.Logger, resp *http.Response) error {
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return fmt.Errorf("Failed to dump response: '%v'", err)
	}
	logPrefix := ">>>"
	logger.Debugf("Received response: (Messages with %s prefix are part of this log)", logPrefix)
	logFriendlyHTTPMessage(logger, string(respDump), logPrefix)

	return nil
}

func logCurl(logger *logger.Logger, req *http.Request) {
	logger.Infof("You can use the following curl command to test this locally")
	logger.Infof("$ %s", httpRequestToCurlString(req))
}

func executeHTTPRequestWithLogging(client *http.Client, req *http.Request, logger *logger.Logger) (*http.Response, error) {
	err := dumpRequest(logger, req)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logFriendlyError(logger, err)
		return nil, fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}
	err = dumpResponseWithBody(logger, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func requestResponseWithoutBody(client *http.Client, url string, statusCode int, logger *logger.Logger) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	err = dumpRequest(logger, req)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}
	err = dumpResponse(logger, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != statusCode {
		return fmt.Errorf("Expected status code %d, got %d", statusCode, resp.StatusCode)
	}
	return nil
}

func validateContent(resp http.Response, expectedContent string, expectedContentType string) error {
	contentLength := len(expectedContent)

	// test content-type

	receivedContentType := resp.Header.Get("Content-Type")
	if receivedContentType == "" {
		return fmt.Errorf("Content-Type header not present")
	}

	if receivedContentType != expectedContentType {
		return fmt.Errorf("Expected content type %q, got %q", expectedContentType, receivedContentType)
	}

	// test content-length

	receivedContentLength := resp.Header.Get("Content-Length")
	if receivedContentLength == "" {
		return fmt.Errorf("Content-Length header not present")
	}

	if receivedContentLength != fmt.Sprintf("%d", contentLength) {
		return fmt.Errorf("Expected content length %d, got %s", contentLength, receivedContentLength)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != expectedContent {
		return fmt.Errorf("Expected the content to be %q got %q", expectedContent, body)
	}

	return nil
}
