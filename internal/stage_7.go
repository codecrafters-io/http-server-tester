package internal

import (
	"fmt"
	"net/http"
	"os"

	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testGetFile(stageHarness *test_case_harness.TestCaseHarness) error {
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
	err := createFileWith(DATA_DIR+fileName, fileContent)
	defer os.Remove(DATA_DIR + fileName)

	if err != nil {
		return err
	}

	err = testGetFileResponse(logger, fileName, fileContent)
	if err != nil {
		return err
	}

	logger.Infof("Testing non existent file returns 404")
	nonExistentFileName := randomFileNameWithPrefix("non-existent")
	err = testNonExistentFileResponseIs404(logger, nonExistentFileName)
	if err != nil {
		return err
	}

	return nil
}

func testNonExistentFileResponseIs404(logger *logger.Logger, fileName string) error {
	httpClient := NewHTTPClient()

	req, err := http.NewRequest("GET", URL+"files/"+fileName, nil)
	if err != nil {
		return err
	}

	resp, err := executeHTTPRequestWithLogging(httpClient, req, logger)
	if err != nil {
		return err
	}

	if resp.StatusCode != 404 {
		return fmt.Errorf("Expected status code 404, got %d", resp.StatusCode)
	}

	return nil
}
