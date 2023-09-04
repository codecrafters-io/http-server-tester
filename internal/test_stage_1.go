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
	logger.Infof("Running stage 1")

	logger.Infof("Sending Get request")
	client := http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:4221", nil)
	if err != nil {
		logFriendlyError(stageHarness.Logger, err)
		return fmt.Errorf("Error creating request:", err)
	}
	req.Close = true

	// Send the request
	response, err := client.Do(req)
	if err != nil {
		logFriendlyError(stageHarness.Logger, err)
		return fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Expected status code 200, got %d", response.StatusCode)
	}

	logger.Debugf("Connection successful")

	return nil
}
