package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"

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
	userAgent := randSeq(20)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Could not create request: %v", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Could not send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Could not read response body: %v", err)
	}

	responseBody := string(body)

	if userAgent != strings.TrimSpace(responseBody) {
		return fmt.Errorf("Custom User-Agent '%s' not found in the response body.\nResponse body: %s", userAgent, responseBody)
	}

	return nil
}
