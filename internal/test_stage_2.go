package internal

import (
	"fmt"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func test404NotFound(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	httpClient := NewHTTPClient()

	response, err := httpClient.Get("http://localhost:4221/random")
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}

	if response.StatusCode != 404 {
		return fmt.Errorf("Expected status code 404, got %d", response.StatusCode)
	}

	logger.Debugf("Connection successful")

	return nil
}
