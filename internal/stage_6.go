package internal

import (
	"github.com/codecrafters-io/tester-utils/random"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHandlesMultipleConcurrentConnections(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	connectionCount := random.RandomInt(2, 3)

	requestResponsePair, err := GetBaseURLGetRequestResponsePair()
	if err != nil {
		return err
	}

	testCase := test_cases.SendRequestTestCase{
		Request:                   requestResponsePair.Request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(*requestResponsePair.Response),
		ShouldSkipUnreadDataCheck: false,
	}

	connections, err := spawnConnections(stageHarness, connectionCount, logger)
	if err != nil {
		return err
	}
	logger.Debugf("Sending first set of requests")
	for i := connectionCount - 1; i >= 0; i-- {
		// Test connections in reverse order so that we don't accidentally test the listen backlog
		// Ref: https://github.com/codecrafters-io/http-server-tester/pull/60
		if err := testCase.RunWithConn(connections[i], logger); err != nil {
			return err
		}

		logger.Debugf("Closing connection %d", i+1)
		err := connections[i].Close()
		if err != nil {
			logFriendlyError(logger, err)
			return err
		}
	}

	// At this point, we have closed all open connections.
	// But the server should still be running.
	// We will now spawn new connections and send requests again.
	connections, err = spawnConnections(stageHarness, connectionCount, logger)
	if err != nil {
		return err
	}
	logger.Debugf("Sending second set of requests")
	for i := range connectionCount {
		if err := testCase.RunWithConn(connections[i], logger); err != nil {
			return err
		}

		logger.Debugf("Closing connection %d", i+1)
		err := connections[i].Close()
		if err != nil {
			logFriendlyError(logger, err)
			return err
		}
	}

	return nil
}
