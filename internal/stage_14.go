package internal

import (
	"fmt"
	"net/http"
	"os"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPersistence3(stageHarness *test_case_harness.TestCaseHarness) error {
	setupDataDirectory()
	defer os.RemoveAll(DATA_DIR)
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run("--directory", DATA_DIR); err != nil {
		return err
	}

	logger := stageHarness.Logger

	// N connections & N unique requests
	connectionCount := 2
	uniqueRequestCount := connectionCount

	requestResponsePairs, err := GetRandomRequestResponsePairs(uniqueRequestCount, logger)
	if err != nil {
		return err
	}

	testCases := make([]test_cases.SendRequestTestCase, uniqueRequestCount)
	requests := make([]*http.Request, uniqueRequestCount)

	for i, requestResponsePair := range requestResponsePairs {
		response := *requestResponsePair.Response

		if i == 1 {
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

	connections, err := spawnConnections(stageHarness, 2, logger)
	if err != nil {
		return err
	}

	logger.Debugf("Sending first set of requests to connection #0")
	for i := range uniqueRequestCount {
		if err := testCases[i].RunWithConn(connections[0], logger); err != nil {
			return err
		}
	}

	if connections[0].IsOpen() {
		return fmt.Errorf("connection is still open")
	}

	logger.Debugf("Sending second set of requests to connection #1")
	for i := range uniqueRequestCount {
		if err := testCases[i].RunWithConn(connections[1], logger); err != nil {
			return err
		}
	}
	if connections[1].IsOpen() {
		return fmt.Errorf("connection is still open")
	}

	// Connections should be closed by the time we get here
	return nil
}
