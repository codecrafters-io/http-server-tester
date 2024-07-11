package internal

import (
	"fmt"
	"net/http"

	"github.com/codecrafters-io/tester-utils/logger"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHandlesMultipleConcurrentConnections(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	connectionCount := random.RandomInt(2, 3)

	request, _ := http.NewRequest("GET", URL, nil)
	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine}
	testCase := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}

	connections, err := spawnConnections(stageHarness, connectionCount, logger)
	if err != nil {
		return err
	}
	logger.Debugf("Sending first set of requests")
	for i := connectionCount; i < 0; i-- {
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
	for i := 0; i < connectionCount; i++ {
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

func spawnConnections(stageHarness *test_case_harness.TestCaseHarness, connectionCount int, logger *logger.Logger) ([]*http_connection.HttpConnection, error) {
	logger.Infof("Creating %d parallel connections", connectionCount)
	connections := make([]*http_connection.HttpConnection, connectionCount)

	for i := 0; i < connectionCount; i++ {
		logger.Debugf("Creating connection %d", i+1)
		conn, err := http_connection.NewInstrumentedHttpConnection(stageHarness, TCP_DEST, fmt.Sprintf("client-%d", i+1))
		if err != nil {
			logFriendlyError(logger, err)
			return nil, err
		}
		connections[i] = conn
	}
	return connections, nil
}
