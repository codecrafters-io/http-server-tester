package internal

import (
	"fmt"
	testerutils "github.com/codecrafters-io/tester-utils"
)

func test200OK(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	httpClient := NewHTTPClient()
	response, err := httpClient.Get("http://localhost:4221")
	if err != nil {
		logFriendlyError(stageHarness.Logger, err)
		return fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Expected status code 200, got %d", response.StatusCode)
	}

	logger := stageHarness.Logger
	logger.Debugf("Connection successful")

	return nil
}
