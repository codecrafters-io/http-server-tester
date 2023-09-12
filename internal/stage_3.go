package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
)

const randomUrlLength = 20

func test404NotFound(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	httpClient := NewHTTPClient()

	var url = URL + randomUrlPath()

	return requestWithStatus(httpClient, url, 404, logger)
}
