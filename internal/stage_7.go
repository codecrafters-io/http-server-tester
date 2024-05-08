package internal

import (
	"fmt"
	"net/http"
	"os"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testGetFile(stageHarness *test_case_harness.TestCaseHarness) error {
	err := os.MkdirAll(DATA_DIR, 0755)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(DATA_DIR)

	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run("--directory", DATA_DIR); err != nil {
		return err
	}

	logger := stageHarness.Logger

	fileName := randomFileName()
	fileContent := randomFileContent()

	logger.Infof("Testing existing file")
	logger.Debugf("Creating file %s in %s", fileName, DATA_DIR)
	logger.Debugf("File Content: %q", fileContent)
	err = createFileWith(DATA_DIR+fileName, fileContent)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("GET", URL+"files/"+fileName, nil)
	if err != nil {
		return err
	}
	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 200, Reason: "OK"}
	header1 := http_parser.Header{Key: "Content-Type", Value: "application/octet-stream"}
	header2 := http_parser.Header{Key: "Content-Length", Value: fmt.Sprintf("%d", len(fileContent))}
	expectedHeaders := []http_parser.Header{header1, header2}
	expectedBody := []byte(fileContent)
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine, Headers: expectedHeaders, Body: expectedBody}

	test_case := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	if err := test_case.Run(stageHarness, TCP_DEST, logger); err != nil {
		return err
	}

	logger.Successf("First test passed.")
	logger.Infof("Testing non existent file returns 404")

	nonExistentFileName := randomFileNameWithPrefix("non-existent")
	request, err = http.NewRequest("GET", URL+"files/"+nonExistentFileName, nil)
	if err != nil {
		return err
	}
	expectedStatusLine = http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 404, Reason: "Not Found"}
	expectedResponse = http_parser.HTTPResponse{StatusLine: expectedStatusLine}

	test_case = test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}
	return test_case.Run(stageHarness, TCP_DEST, logger)
}
