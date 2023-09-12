package internal

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"

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

func sendRequest(client *http.Client, url string, logger *testerutils.Logger) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Sending Request:\n%s", string(reqDump))

	resp, err := client.Do(req)
	if err != nil {
		logFriendlyError(logger, err)
		return nil, fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Received Response:\n%s", string(respDump))
	return resp, nil
}

func requestWithStatus(client *http.Client, url string, statusCode int, logger *testerutils.Logger) error {
	resp, err := sendRequest(client, url, logger)
	if err != nil {
		return err
	}

	if resp.StatusCode != statusCode {
		return fmt.Errorf("Expected status code %d, got %d", statusCode, resp.StatusCode)
	}
	return nil
}
