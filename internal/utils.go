package internal

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/logger"

	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func spawnConnection(stageHarness *test_case_harness.TestCaseHarness, logger *logger.Logger) (*http_connection.HttpConnection, error) {
	logger.Debugf("Creating connection")

	conn, err := http_connection.NewInstrumentedHttpConnection(stageHarness, TCP_DEST, "")
	if err != nil {
		logFriendlyError(logger, err)
		return nil, err
	}

	return conn, nil
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

func repeatSingleConnection(stageHarness *test_case_harness.TestCaseHarness, connectionCount int, logger *logger.Logger) ([]*http_connection.HttpConnection, error) {
	logger.Infof("Repeating single connection %d times", connectionCount)
	connections := make([]*http_connection.HttpConnection, connectionCount)

	conn, err := http_connection.NewInstrumentedHttpConnection(stageHarness, TCP_DEST, "client")
	if err != nil {
		logFriendlyError(logger, err)
		return nil, err
	}

	for i := 0; i < connectionCount; i++ {
		logger.Debugf("Repeating connection %d", i+1)
		connections[i] = conn
	}

	return connections, nil
}

func assertConnectionIsOpen(connection *http_connection.HttpConnection, logger *logger.Logger) error {
	if !connection.IsOpen() {
		return fmt.Errorf("Expected connection to stay open, but it is closed")
	}
	return nil
}
