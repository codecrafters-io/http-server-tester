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
	logger.Infof("Running stage 3")

	httpClient := NewHTTPClient()

	// There are no urls that would be this long so there is 0 probability of collison ever
	var url = URL + randSeq(randomUrlLength)

	logger.Infof("Calling %s", url)
	requestWithStatus(httpClient, url, 404, logger)

	return nil
}
