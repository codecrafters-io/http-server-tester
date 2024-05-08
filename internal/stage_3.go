package internal

import (
	"net/http"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func test404NotFound(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	var url = URL + randomUrlPath()

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 404, Reason: "Not Found"}
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine}

	test_case := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	return test_case.Run(stageHarness, TCP_DEST, logger)
}
