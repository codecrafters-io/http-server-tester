package internal

import (
	"fmt"
	"net/http"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testRespondWithUserAgent(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewHTTPClient()

	url := URL + "user-agent"
	userAgent := randomUserAgent()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Could not create request: %v", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := executeHTTPRequestWithLogging(client, req, logger)
	if err != nil {
		return err
	}

	err = validateContent(*resp, userAgent, "text/plain")
	if err != nil {
		return err
	}

	return nil
}
