package test_cases

import (
	"fmt"
	"net/http"
	"strings"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type SendRequestTestCase struct {
	Request                   *http.Request
	Assertion                 http_assertions.HTTPResponseAssertion
	ShouldSkipUnreadDataCheck bool

	// ReceivedResponse is set after the test case is run
	ReceivedResponse http_parser.HTTPResponse
}

func (t *SendRequestTestCase) Run(stageHarness *test_case_harness.TestCaseHarness, address string, logger *logger.Logger) error {
	conn, err := http_connection.NewInstrumentedHttpConnection(stageHarness, address, "")
	if err != nil {
		return fmt.Errorf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	addr := strings.Split(address, ":")
	host, port := addr[0], addr[1]
	logger.Debugf(fmt.Sprintf("Connected to %s port %s", host, port))

	err = conn.SendRequest(t.Request)
	if err != nil {
		return fmt.Errorf("Failed to send request: %v", err)
	}

	response, err := conn.ReadResponse()
	if err != nil {
		return fmt.Errorf("Failed to read response: \n%v", err)
	}
	t.ReceivedResponse = response

	if err = t.Assertion.Run(response, logger); err != nil {
		return err
	}

	if !t.ShouldSkipUnreadDataCheck {
		conn.EnsureNoUnreadData()
	}

	return nil
}

func (t *SendRequestTestCase) RunWithConn(conn *http_connection.HttpConnection, logger *logger.Logger) error {
	err := conn.SendRequest(t.Request)
	if err != nil {
		return fmt.Errorf("Failed to send request: %v", err)
	}

	response, err := conn.ReadResponse()
	if err != nil {
		return fmt.Errorf("Failed to read response: \n%v", err)
	}
	t.ReceivedResponse = response

	if err = t.Assertion.Run(response, logger); err != nil {
		return err
	}

	if !t.ShouldSkipUnreadDataCheck {
		conn.EnsureNoUnreadData()
	}

	return nil
}
