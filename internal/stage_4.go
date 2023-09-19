package internal

import (
	"fmt"
	"net/http"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testRespondWithContent(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	httpClient := NewHTTPClient()

	content := randomUrlPath()
	url := URL + "echo/" + content

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := sendRequest(httpClient, req, logger)

	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	if err := validateContent(*resp, content, "text/plain"); err != nil {
		return err
	}

	return nil
}
