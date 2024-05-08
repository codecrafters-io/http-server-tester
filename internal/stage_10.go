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

func testRespondWithCorrectContentEncoding(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	conn1, err := http_connection.NewInstrumentedHttpConnection(stageHarness, TCP_DEST, "client")
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to create connection: %v", err)
	}
	defer conn1.Close()
	logger.Debugln("Connection established, sending request...")

	content := randomUrlPath()
	url := URL + "echo/" + content

	// Success case : 200 OK with Content-Encoding: gzip
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Could not create request: %v", err)
	}
	request.Header.Set("Accept-Encoding", "encoding-1, gzip, encoding-2")

	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	header := http_parser.Header{Key: "Content-Encoding", Value: "gzip"}
	headerFormattedAsString := fmt.Sprintf("%s: %s", header.Key, header.Value)
	expectedHeaders := []http_parser.Header{header}
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine, Headers: expectedHeaders}

	test_case := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	if err := test_case.Run(conn1, logger, " "+headerFormattedAsString); err != nil {
		return err
	}

	conn2, err := http_connection.NewInstrumentedHttpConnection(stageHarness, TCP_DEST, "client")
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to create connection: %v", err)
	}
	defer conn2.Close()
	logger.Debugln("Connection established, sending request...")

	// Failure case : 200 OK without Content-Encoding: gzip
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Could not create request: %v", err)
	}
	request.Header.Set("Accept-Encoding", "encoding-1, encoding-2")

	expectedStatusLine = http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	expectedResponse = http_parser.HTTPResponse{StatusLine: expectedStatusLine}

	test_case = test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	if err := test_case.Run(conn2, logger, " "); err != nil {
		return err
	}

	if test_case.ReceivedResponse.FindHeader("Content-Encoding") != "" {
		return fmt.Errorf("Content-Encoding header should not be present")
	}

	return nil
}
