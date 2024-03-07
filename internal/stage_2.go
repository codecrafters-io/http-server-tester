package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func test200OK(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	httpClient := NewHTTPClient()

	return requestResponseWithoutBody(httpClient, URL, 200, logger)
}
