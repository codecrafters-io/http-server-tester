package internal

import (
	"fmt"
	"net/http"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRespondWithContentEncoding(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	content := randomUrlPath()
	url := URL + "echo/" + content

	// Success case : 200 OK with Content-Encoding: gzip
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Could not create request: %v", err)
	}
	request.Header.Set("Accept-Encoding", "gzip")

	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	header := http_parser.Header{Key: "Content-Encoding", Value: "gzip"}
	expectedHeaders := []http_parser.Header{header}
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine, Headers: expectedHeaders}

	test_case := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	if err := test_case.Run(stageHarness, TCP_DEST, logger); err != nil {
		return err
	}
	logger.Successf("First test passed.")

	// Failure case : 200 OK without Content-Encoding: gzip
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Could not create request: %v", err)
	}
	request.Header.Set("Accept-Encoding", "invalid-encoding")

	expectedStatusLine = http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	expectedResponse = http_parser.HTTPResponse{StatusLine: expectedStatusLine}

	test_case = test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	if err := test_case.Run(stageHarness, TCP_DEST, logger); err != nil {
		return err
	}

	if test_case.ReceivedResponse.FindHeader("Content-Encoding") != "" {
		return fmt.Errorf("Content-Encoding header should not be present")
	}
	logger.Successf("✓ Content-Encoding header is not present")

	return nil
}
