package internal

import (
	"net/http"
	"os"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

// TODO: Add better emulated curl logs
func testPersistence1(stageHarness *test_case_harness.TestCaseHarness) error {
	setupDataDirectory()
	defer os.RemoveAll(DATA_DIR)
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run("--directory", DATA_DIR); err != nil {
		return err
	}

	logger := stageHarness.Logger

	uniqueRequestCount := 2
	requestResponsePairs, err := GetRandomRequestResponsePairs(uniqueRequestCount, logger)
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

	connection, err := spawnPersistentConnection(stageHarness, logger)
	if err != nil {
		return err
	}

	logger.Debugf("Sending first set of requests")
	logger.Infof("$ %s", http_connection.HttpKeepAliveRequestToCurlStringForMultipleRequests(requests))
	for _, testCase := range testCases {
		if err := testCase.RunWithConn(connection, logger); err != nil {
			return err
		}
	}

	logger.Debugf("Sending second set of requests")
	logger.Infof("$ %s", http_connection.HttpKeepAliveRequestToCurlStringForMultipleRequests(requests))
	for _, testCase := range testCases {
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
