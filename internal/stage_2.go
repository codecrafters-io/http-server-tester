package internal

import (
	"fmt"
	"net/http"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func test200OK(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	httpClient := NewHTTPClient()

	requestWithStatus(httpClient, URL, 200, logger)

	return nil
}

func requestWithStatus(client *http.Client, url string, statusCode int, logger *testerutils.Logger) error {
	response, err := client.Get(url)
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}
	if response.StatusCode != statusCode {
		return fmt.Errorf("Expected status code %d, got %d", statusCode, response.StatusCode)
	}
	return nil
}
