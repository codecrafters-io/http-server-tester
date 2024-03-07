package internal

import (
	"fmt"

	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func antiCheatBasic(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	client := NewHTTPClient()

	resp, err := client.Get(URL)
	if err != nil {
		return nil
	}

	if resp.Proto != "HTTP/1.1" {
		return fail(logger)
	}

	if date := resp.Header.Get("Date"); date != "" {
		return fail(logger)
	}

	if server := resp.Header.Get("Server"); server != "" {
		return fail(logger)
	}

	return nil
}

func fail(logger *logger.Logger) error {
	logger.Criticalf("anti-cheat (ac1) failed.")
	logger.Criticalf("Are you sure you aren't running this against an actual HTTP server?")
	return fmt.Errorf("anti-cheat (ac1) failed")
}
