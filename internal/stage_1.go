package internal

import (
	"fmt"
	"net"
	"time"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testConnects(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	// Friendly logs for the first stage - this doesn't have to be done for further stages
	var conn net.Conn
	retries := 0
	var err error
	for {
		logger.Infof("Connecting to %s using TCP", TCP_DEST)
		conn, err = net.Dial("tcp", TCP_DEST)
		if err != nil && retries > 15 {
			logger.Infof("All retries failed.")
			return err
		}

		if err != nil {
			if b.HasExited() {
				return fmt.Errorf("Looks like your program has terminated. A HTTP server is expected to be a long-running process.")
			}

			// Don't print errors in the first second
			if retries > 2 {
				logger.Infof("Failed to connect to port 4221, retrying in 1s")
			}

			retries += 1
			time.Sleep(1000 * time.Millisecond)
		} else {
			logger.Infof("Success! Closing connection")
			conn.Close()
			break
		}
	}

	return nil
}
