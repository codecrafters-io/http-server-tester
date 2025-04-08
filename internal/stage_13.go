package internal

import (
	"net/http"

	"github.com/codecrafters-io/tester-utils/random"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPersistence2(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	// N connections & N unique requests
	connectionCount := random.RandomInt(2, 3)
	uniqueRequestCount := connectionCount

	requestResponsePairs, err := GetRandomRequestResponsePairs(uniqueRequestCount)
	if err != nil {
		return err
	}
	testCases := make([]test_cases.SendRequestTestCase, uniqueRequestCount)
	requests := make([]*http.Request, uniqueRequestCount)

	for i, requestResponsePair := range requestResponsePairs {
		testCases[i] = test_cases.SendRequestTestCase{
			Request:                   requestResponsePair.Request,
			Assertion:                 http_assertions.NewHTTPResponseAssertion(*requestResponsePair.Response),
			ShouldSkipUnreadDataCheck: false,
		}
		requests[i] = requestResponsePair.Request
	}

	connections, err := spawnPersistentConnections(stageHarness, connectionCount, logger)
	if err != nil {
		return err
	}

	logger.Debugf("Sending first set of requests")
	logger.Infof("$ %s", http_connection.HttpKeepAliveRequestToCurlString(requests))
	for i := connectionCount - 1; i >= 0; i-- {
		// Test connections in reverse order so that we don't accidentally test the listen backlog
		// Ref: https://github.com/codecrafters-io/http-server-tester/pull/60
		if err := testCases[i].RunWithConn(connections[i], logger); err != nil {
			return err
		}
	}

	logger.Debugf("Sending second set of requests")
	logger.Infof("$ %s", http_connection.HttpKeepAliveRequestToCurlString(requests))
	for i := range connectionCount {
		logger.Debugf("* Re-using existing connection with host localhost")
		if err := testCases[i].RunWithConn(connections[i], logger); err != nil {
			return err
		}
	}

	for _, connection := range connections {
		err := connection.Close()
		if err != nil {
			logFriendlyError(logger, err)
			return err
		}
	}

	return nil
}
