package internal

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func test200OK(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	httpClient := NewHTTPClient()

	return requestWithStatus(httpClient, URL, 200, logger)
}

func requestWithStatus(client *http.Client, url string, statusCode int, logger *testerutils.Logger) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err
	}
	logger.Debugf("Sending Request:\n%s", string(reqDump))

	resp, err := client.Do(req)
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	logger.Debugf("Received Response:\n%s", string(respDump))

	if resp.StatusCode != statusCode {
		return fmt.Errorf("Expected status code %d, got %d", statusCode, resp.StatusCode)
	}
	return nil
}
