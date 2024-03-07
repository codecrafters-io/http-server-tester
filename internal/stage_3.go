package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

const randomUrlLength = 20

func test404NotFound(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	httpClient := NewHTTPClient()

	var url = URL + randomUrlPath()

	return requestResponseWithoutBody(httpClient, url, 404, logger)
}
