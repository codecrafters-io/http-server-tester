package internal

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPostFile(stageHarness *test_case_harness.TestCaseHarness) error {
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

	err = postFile(logger, fileName, fileContent)
	if err != nil {
		return err
	}

	err = validateFile(logger, fileName, fileContent)
	if err != nil {
		return err
	}

	return nil
}

func postFile(logger *logger.Logger, fileName string, fileContent string) error {
	httpClient := NewHTTPClient()

	req, err := http.NewRequest("POST", URL+"files/"+fileName, bytes.NewBufferString(fileContent))
	if err != nil {
		return err
	}
	err = dumpRequest(logger, req)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}
	err = dumpResponse(logger, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("Expected status code 201, got %d", resp.StatusCode)
	}

	return nil
}
