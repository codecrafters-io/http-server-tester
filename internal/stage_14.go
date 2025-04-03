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

	request, _ := http.NewRequest("GET", URL, nil)
	request.Header.Set("Connection", "close")

	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	expectedHeaders := []http_parser.Header{{Key: "Connection", Value: "close"}}
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine, Headers: expectedHeaders}
	testCase := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}

	// TODO: We don't want to same the request N times
	// Need to implement random request generator
	requestCount := 1
	connections, err := spawnConnections(stageHarness, 2, logger)
	if err != nil {
		return err
	}

	logger.Debugf("Sending first set of requests")
	for range requestCount {
		if err := testCase.RunWithConn(connections[0], logger); err != nil {
			return err
		}
	}
	if connections[0].IsOpen() {
		return fmt.Errorf("connection is still open")
	}

	logger.Debugf("Sending second set of requests")
	for range requestCount {
		if err := testCase.RunWithConn(connections[1], logger); err != nil {
			return err
		}
	}
	if connections[1].IsOpen() {
		return fmt.Errorf("connection is still open")
	}

	for _, connection := range connections {
		err = connection.Close()
		if err != nil {
			logFriendlyError(logger, err)
			return err
		}
	}

	return nil
}
