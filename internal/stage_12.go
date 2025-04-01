package internal

import (
	"net/http"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPersistence1(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	request, _ := http.NewRequest("GET", URL, nil)
	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine}
	testCase := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}

	connections, err := repeatSingleConnection(stageHarness, 2, logger)
	if err != nil {
		return err
	}

	logger.Debugf("Sending first set of requests")
	for i := range connections {
		// TODO: Reverse all loops
		// Test connections in reverse order so that we don't accidentally test the listen backlog
		// Ref: https://github.com/codecrafters-io/http-server-tester/pull/60
		if err := testCase.RunWithConn(connections[i], logger); err != nil {
			return err
		}
	}

	logger.Debugf("Sending second set of requests")
	for i := range connections {
		// Test connections in reverse order so that we don't accidentally test the listen backlog
		// Ref: https://github.com/codecrafters-io/http-server-tester/pull/60
		if err := testCase.RunWithConn(connections[i], logger); err != nil {
			return err
		}
	}

	err = connections[0].Close()
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	return nil
}
