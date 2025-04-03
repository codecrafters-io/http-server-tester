package internal

import (
	"net/http"
	"os"

	"github.com/codecrafters-io/tester-utils/random"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPersistence2(stageHarness *test_case_harness.TestCaseHarness) error {
	setupDataDirectory()
	defer os.RemoveAll(DATA_DIR)
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run("--directory", DATA_DIR); err != nil {
		return err
	}

	logger := stageHarness.Logger
	connectionCount := random.RandomInt(2, 3)
	requestCount := 2

	request, _ := http.NewRequest("GET", URL, nil)
	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine}
	testCase := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	connections, err := spawnPersistentConnections(stageHarness, connectionCount, logger)
	if err != nil {
		return err
	}

	logger.Debugf("Sending first set of requests")
	logger.Infof("$ %s", http_connection.HttpKeepAliveRequestToCurlString(request, requestCount))
	for i := connectionCount - 1; i >= 0; i-- {
		// Test connections in reverse order so that we don't accidentally test the listen backlog
		// Ref: https://github.com/codecrafters-io/http-server-tester/pull/60
		for range requestCount {
			if err := testCase.RunWithConn(connections[i], logger); err != nil {
				return err
			}
		}
	}

	logger.Debugf("Sending second set of requests")
	logger.Infof("$ %s", http_connection.HttpKeepAliveRequestToCurlString(request, requestCount))
	for i := connectionCount - 1; i >= 0; i-- {
		// Test connections in reverse order so that we don't accidentally test the listen backlog
		// Ref: https://github.com/codecrafters-io/http-server-tester/pull/60
		for range requestCount {
			if err := testCase.RunWithConn(connections[i], logger); err != nil {
				return err
			}
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
