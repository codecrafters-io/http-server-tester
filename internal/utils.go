package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/tester-utils/logger"

	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func setupDataDirectory() {
	err := os.MkdirAll(DATA_DIR, 0755)
	if err != nil {
		panic(err)
	}
}

func spawnPersistentConnection(stageHarness *test_case_harness.TestCaseHarness, logger *logger.Logger) (*http_connection.HttpConnection, error) {
	logger.Debugf("Creating connection")

	conn, err := http_connection.NewInstrumentedHttpConnection(stageHarness, TCP_DEST, "")

	// We want to log all the requests at once, not one by one
	// Note: No support for logPrefix here
	conn.Callbacks.BeforeSendRequest = nil
	conn.Callbacks.AfterReadResponse = func(response http_parser.HTTPResponse) {
		for _, line := range strings.Split(strings.TrimSpace(response.FormattedString()), "\r\n") {
			stageHarness.Logger.Debugf("%s%s", "< ", line)
		}
		stageHarness.Logger.Debugf("%s%s", "< ", "")
		stageHarness.Logger.Debugf("* Connection #0 to host localhost left intact")
	}

	if err != nil {
		logFriendlyError(logger, err)
		return nil, err
	}

	return conn, nil
}

func spawnPersistentConnections(stageHarness *test_case_harness.TestCaseHarness, connectionCount int, logger *logger.Logger) ([]*http_connection.HttpConnection, error) {
	logger.Debugf("Creating %d persistent connections", connectionCount)
	connections := make([]*http_connection.HttpConnection, connectionCount)

	for i := range connectionCount {
		conn, err := http_connection.NewInstrumentedHttpConnection(stageHarness, TCP_DEST, fmt.Sprintf("client-%d", i+1))
		if err != nil {
			logFriendlyError(logger, err)
			return nil, err
		}

		// We want to log all the requests at once, not one by one
		// Note: No support for logPrefix here
		conn.Callbacks.BeforeSendRequest = nil
		conn.Callbacks.AfterReadResponse = func(response http_parser.HTTPResponse) {
			for _, line := range strings.Split(strings.TrimSpace(response.FormattedString()), "\r\n") {
				stageHarness.Logger.Debugf("%s%s", "< ", line)
			}
			stageHarness.Logger.Debugf("%s%s", "< ", "")
			stageHarness.Logger.Debugf("* Connection #%d to host localhost left intact", i)
		}

		connections[i] = conn
	}

	return connections, nil
}

func spawnConnections(stageHarness *test_case_harness.TestCaseHarness, connectionCount int, logger *logger.Logger) ([]*http_connection.HttpConnection, error) {
	logger.Infof("Creating %d parallel connections", connectionCount)
	connections := make([]*http_connection.HttpConnection, connectionCount)

	for i := range connectionCount {
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
