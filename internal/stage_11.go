package internal

import (
	"fmt"
	"net/http"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRespondWithEncodedData(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	content := randomUrlPath()
	url := URL + "echo/" + content

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

	if test_case.ReceivedResponse.FindHeader("Content-Length") != fmt.Sprintf("%d", len(test_case.ReceivedResponse.Body)) {
		return fmt.Errorf("Content-Length header (%v bytes) does not match the length of the body (%d bytes)", test_case.ReceivedResponse.FindHeader("Content-Length"), len(test_case.ReceivedResponse.Body))
	}
	logger.Successf("✓ Content-Length header is present")

	gzipString := test_case.ReceivedResponse.Body
	decodedString, err := decodeGZIP(gzipString)
	if err != nil {
		return fmt.Errorf("Failed to decode gzip: %v", err)
	}
	logger.Successf("✓ Body is gzip encoded")

	if string(decodedString) != content {
		return fmt.Errorf("Expected %s, got %s", content, decodedString)
	}
	logger.Successf("✓ Body is correct")

	return nil
}
