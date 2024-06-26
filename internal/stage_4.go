package internal

import (
	"fmt"
	"net/http"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRespondWithContent(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	content := randomUrlPath()
	url := URL + "echo/" + content

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	header1 := http_parser.Header{Key: "Content-Type", Value: "text/plain"}
	header2 := http_parser.Header{Key: "Content-Length", Value: fmt.Sprintf("%d", len(content))}
	expectedHeaders := []http_parser.Header{header1, header2}
	expectedBody := []byte(content)
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine, Headers: expectedHeaders, Body: expectedBody}

	test_case := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	return test_case.Run(stageHarness, TCP_DEST, logger)
}
