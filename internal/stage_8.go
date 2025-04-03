package internal

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPostFile(stageHarness *test_case_harness.TestCaseHarness) error {
	setupDataDirectory()
	defer os.RemoveAll(DATA_DIR)
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run("--directory", DATA_DIR); err != nil {
		return err
	}

	logger := stageHarness.Logger

	fileName := randomFileName()
	fileContent := randomFileContent()

	request, err := http.NewRequest("POST", URL+"files/"+fileName, bytes.NewBufferString(fileContent))
	request.Header.Add("Content-Length", fmt.Sprint(len(fileContent)))
	request.Header.Add("Content-Type", "application/octet-stream")
	if err != nil {
		return err
	}
	expectedStatusLine := http_parser.StatusLine{Version: "HTTP/1.1", StatusCode: 201, Reason: "Created"}
	expectedResponse := http_parser.HTTPResponse{StatusLine: expectedStatusLine}

	test_case := test_cases.SendRequestTestCase{
		Request:                   request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(expectedResponse),
		ShouldSkipUnreadDataCheck: false,
	}

	if err := test_case.Run(stageHarness, TCP_DEST, logger); err != nil {
		return err
	}

	err = validateFile(logger, fileName, fileContent)
	if err != nil {
		return err
	}

	return nil
}
