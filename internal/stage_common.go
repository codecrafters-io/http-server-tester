package internal

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"

	testerutils "github.com/codecrafters-io/tester-utils"
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

func sendRequest(client *http.Client, req *http.Request, logger *testerutils.Logger) (*http.Response, error) {
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	logger.Infof("Sending request (status line): %s", getFirstLine(string(reqDump)))
	logger.Debugf("Sending request:\n%s", string(reqDump))

	resp, err := client.Do(req)
	if err != nil {
		logFriendlyError(logger, err)
		return nil, fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	logger.Infof("Received response: (status line) %s", getFirstLine(string(respDump)))
	logger.Debugf("Received response:\n%s", string(respDump))
	return resp, nil
}

func requestWithStatus(client *http.Client, url string, statusCode int, logger *testerutils.Logger) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := sendRequest(client, req, logger)
	if err != nil {
		return err
	}

	if resp.StatusCode != statusCode {
		return fmt.Errorf("Expected status code %d, got %d", statusCode, resp.StatusCode)
	}
	return nil
}

func validateContent(resp http.Response, expected string) error {
	contentLength := len(expected)

	// test content-type

	receivedContentType := resp.Header.Get("Content-Type")
	if receivedContentType == "" {
		return fmt.Errorf("Content-Type header not present")
	}

	if receivedContentType != "text/plain" {
		return fmt.Errorf("Expected content type text/plain, got %s", receivedContentType)
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

	if string(body) != expected {
		return fmt.Errorf("Expected the content to be %s got %s", expected, body)
	}

	return nil
}
