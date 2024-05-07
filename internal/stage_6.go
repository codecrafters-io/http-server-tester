package internal

import (
	"fmt"
	"math/rand"
	"net/http"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHandlesMultipleConcurrentConnections(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	connectionCount := rand.Intn(2) + 3 // Fixtures will fail, unless we use tester-utils/random
	logger.Infof("Creating %d parallel connections", connectionCount)

	conns := make([]*http_connection.HttpConnection, connectionCount)

	for i := 0; i < connectionCount; i++ {
		logger.Debugf("Creating connection %d", i+1)
		conn, err := http_connection.NewInstrumentedHttpConnection(stageHarness, TCP_DEST, fmt.Sprintf("client-%d", i+1))
		if err != nil {
			logFriendlyError(logger, err)
			return err
		}
		conns[i] = conn
	}

	request, _ := http.NewRequest("GET", URL, nil)
	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine}
	test_case := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}

	for i := 0; i < connectionCount; i++ {
		if err := test_case.Run(conns[i], logger); err != nil {
			return err
		}
	}

	for i := 0; i < connectionCount; i++ {
		logger.Debugf("Closing connection %d", i+1)
		err := conns[i].Close()
		if err != nil {
			logFriendlyError(logger, err)
			return err
		}
	}

	return nil
}
