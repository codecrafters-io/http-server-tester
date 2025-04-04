package internal

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func antiCheatBasic(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	programLogger := stageHarness.Logger

	client := NewHTTPClient()

	resp, err := client.Get(URL)
	if err != nil {
		return nil
	}

	if resp.Proto != "HTTP/1.1" {
		return fail(programLogger)
	}

	if date := resp.Header.Get("Date"); date != "" {
		return fail(programLogger)
	}

	if server := resp.Header.Get("Server"); server != "" {
		return fail(programLogger)
	}

	return nil
}

func fail(logger *logger.Logger) error {
	logger.Criticalf("anti-cheat (ac1) failed.")
	logger.Criticalf("Please contact us at hello@codecrafters.io if you think is a mistake.")
	return fmt.Errorf("anti-cheat (ac1) failed")
}
