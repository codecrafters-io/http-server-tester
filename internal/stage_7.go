package internal

import (
	"net/http"
	"os"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testGetFile(stageHarness *test_case_harness.TestCaseHarness) error {
	setupDataDirectory()
	defer os.RemoveAll(DATA_DIR)
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run("--directory", DATA_DIR); err != nil {
		return err
	}

	logger := stageHarness.Logger

	requestResponsePair, err := GetFilesRequestResponsePair(logger)
	if err != nil {
		return err
	}

	test_case := test_cases.SendRequestTestCase{
		Request:                   requestResponsePair.Request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(*requestResponsePair.Response),
		ShouldSkipUnreadDataCheck: false,
	}
	if err := test_case.Run(stageHarness, TCP_DEST, logger); err != nil {
		return err
	}

	logger.Successf("First test passed.")
	logger.Infof("Testing non existent file returns 404")

	nonExistentFileName := randomFileNameWithPrefix("non-existent")
	request, err := http.NewRequest("GET", URL+"files/"+nonExistentFileName, nil)
	if err != nil {
		return err
	}
	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 404, Reason: "Not Found"}
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine}

	test_case = test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	return test_case.Run(stageHarness, TCP_DEST, logger)
}
