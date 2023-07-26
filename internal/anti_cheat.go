package internal

import (
	"fmt"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func antiCheatTest(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	// TODO: Figure out how to implement anti-cheat
	if false {
		stageHarness.Logger.Criticalf("anti-cheat (ac1) failed.")
		stageHarness.Logger.Criticalf("Are you sure you aren't running this against an actual HTTP server?")
		return fmt.Errorf("anti-cheat (ac1) failed")
	}

	return nil
}
