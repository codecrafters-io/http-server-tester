package internal

import (
	"net/http"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPersistence1(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	uniqueRequestCount := 2
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

	connection, err := spawnPersistentConnection(stageHarness, logger)
	if err != nil {
		return err
	}

	logger.Infof("$ %s", http_connection.HttpKeepAliveRequestToCurlString(requests))
	for i, testCase := range testCases {
		if i != 0 {
			logger.Debugf("* Re-using existing connection with host localhost")
		}
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
