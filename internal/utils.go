package internal

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/tester-utils/logger"

	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

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
		// TODO: Update conn count + host addr
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

	for i := 0; i < connectionCount; i++ {
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
			// TODO: Update conn count + host addr
			stageHarness.Logger.Debugf("* Connection #%d to host localhost left intact", i+1)
		}

		connections[i] = conn
	}

	return connections, nil
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

func getConnectionURL(conn net.Conn) string {
	// Get the remote address (server address)
	remoteAddr := conn.RemoteAddr().String()

	// Parse the address to get host and port
	host, port, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr // fallback to full address if parsing fails
	}

	// Construct URL
	return fmt.Sprintf("http://%s:%s", host, port)
}
