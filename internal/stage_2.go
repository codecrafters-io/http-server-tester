package internal

import (
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
