package internal

import (
	"fmt"
	"net/http"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPersistence3(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	// 1 connections & N unique requests
	uniqueRequestCount := 2

	requestResponsePairs, err := GetRandomRequestResponsePairs(uniqueRequestCount)
	if err != nil {
		return err
	}

	testCases := make([]test_cases.SendRequestTestCase, uniqueRequestCount)
	requests := make([]*http.Request, uniqueRequestCount)

	for i, requestResponsePair := range requestResponsePairs {
		response := *requestResponsePair.Response

		if i == uniqueRequestCount-1 {
			// Close connection on the last request
			// Manually add & assert for the Connection header
			requestResponsePair.Request.Header.Set("Connection", "close")
			response.Headers = append(response.Headers, http_parser.Header{Key: "Connection", Value: "close"})
		}

		testCases[i] = test_cases.SendRequestTestCase{
			Request:                   requestResponsePair.Request,
			Assertion:                 http_assertions.NewHTTPResponseAssertion(response),
			ShouldSkipUnreadDataCheck: false,
		}
		requests[i] = requestResponsePair.Request
	}

	connection, err := spawnConnection(stageHarness, logger)
	if err != nil {
		return err
	}

	for i := range uniqueRequestCount {
		if err := testCases[i].RunWithConn(connection, logger); err != nil {
			return err
		}
	}

	if connection.IsOpen() {
		return fmt.Errorf("connection is still open")
	} else {
		logger.Successf("Connection #0 is closed")
	}

	// Connections should be closed by the time we get here
	return nil
}
