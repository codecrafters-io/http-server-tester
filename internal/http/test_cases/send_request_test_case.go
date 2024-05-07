package test_cases

import (
	"fmt"
	"net/http"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/tester-utils/logger"
)

type SendRequestTestCase struct {
	Request                   *http.Request
	Assertion                 http_assertions.HTTPResponseAssertion
	ShouldSkipUnreadDataCheck bool

	// ReceivedResponse is set after the test case is run
	ReceivedResponse http_parser.HTTPResponse
}

func (t *SendRequestTestCase) Run(conn *http_connection.HttpConnection, logger *logger.Logger) error {
	err := conn.SendRequest(t.Request)
	if err != nil {
		return fmt.Errorf("Failed to send request: %v", err)
	}

	response, err := conn.ReadResponse()
	if err != nil {
		return fmt.Errorf("Failed to read response: %v", err)
	}
	t.ReceivedResponse = response

	if err = t.Assertion.Run(response); err != nil {
		return err
	}

	if !t.ShouldSkipUnreadDataCheck {
		conn.EnsureNoUnreadData()
	}

	logger.Successf("Received %s", response.FormattedString())
	return nil
}
