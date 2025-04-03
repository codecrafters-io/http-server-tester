package internal

import (
	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

// TODO: Add better emulated curl logs
func testPersistence1(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	requestResponsePair, err := GetBaseURLGetRequestResponsePair()
	if err != nil {
		return err
	}

	testCase := test_cases.SendRequestTestCase{
		Request:                   requestResponsePair.Request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(*requestResponsePair.Response),
		ShouldSkipUnreadDataCheck: false,
	}

	// TODO: We don't want to same the request N times
	// Need to implement random request generator
	requestCount := 2
	connection, err := spawnPersistentConnection(stageHarness, logger)
	if err != nil {
		return err
	}

	logger.Debugf("Sending first set of requests")
	logger.Infof("$ %s", http_connection.HttpKeepAliveRequestToCurlString(requestResponsePair.Request, requestCount))
	for range requestCount {
		if err := testCase.RunWithConn(connection, logger); err != nil {
			return err
		}
	}

	logger.Debugf("Sending second set of requests")
	logger.Infof("$ %s", http_connection.HttpKeepAliveRequestToCurlString(requestResponsePair.Request, requestCount))
	for range requestCount {
		if err := testCase.RunWithConn(connection, logger); err != nil {
			return err
		}
	}

	err = connection.Close()
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	return nil
}
