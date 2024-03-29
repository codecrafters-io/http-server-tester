package internal

import (
	"fmt"
	"net/http"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRespondWithContent(stageHarness *test_case_harness.TestCaseHarness) error {
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

	resp, err := executeHTTPRequestWithLogging(httpClient, req, logger)

	if err != nil {
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
