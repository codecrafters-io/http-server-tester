package internal

import (
	"fmt"
	"net/http"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_connection "github.com/codecrafters-io/http-server-tester/internal/http/connection"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRespondWithUserAgent(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	conn, err := http_connection.NewInstrumentedHttpConnection(stageHarness, TCP_DEST, "client")
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to create connection: %v", err)
	}
	defer conn.Close()
	logger.Debugln("Connection established, sending request...")

	url := URL + "user-agent"
	userAgent := randomUserAgent()

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Could not create request: %v", err)
	}
	request.Header.Set("User-Agent", userAgent)

	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	header1 := http_parser.Header{Key: "Content-Type", Value: "text/plain"}
	header2 := http_parser.Header{Key: "Content-Length", Value: fmt.Sprintf("%d", len(userAgent))}
	expectedHeaders := []http_parser.Header{header1, header2}
	expectedBody := []byte(userAgent)
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine, Headers: expectedHeaders, Body: expectedBody}

	test_case := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	return test_case.Run(conn, logger, " "+userAgent)
}
