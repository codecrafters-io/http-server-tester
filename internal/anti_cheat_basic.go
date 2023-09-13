package internal

import (
	"fmt"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func antiCheatBasic(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	fmt.Println("Running anti-cheat (ac1).")

	client := NewHTTPClient()

	resp, err := client.Get(URL)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("Response Headers:")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	if resp.Proto != "HTTP/1.1" {
		fmt.Println(resp.Proto)
		return fail(logger)
	}

	if date := resp.Header.Get("Date"); date != "" {
		fmt.Println(date)
		return fail(logger)
	}

	if server := resp.Header.Get("Server"); server != "" {
		return fail(logger)
	}

	return nil
}

func fail(logger *testerutils.Logger) error {
	logger.Criticalf("anti-cheat (ac1) failed.")
	logger.Criticalf("Are you sure you aren't running this against an actual HTTP server?")
	return fmt.Errorf("anti-cheat (ac1) failed")
}
