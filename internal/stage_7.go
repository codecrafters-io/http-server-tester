package internal

import (
	"fmt"
	"net/http"
	"os"

	testerutils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
)

func testGetFile(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run("--directory", DATA_DIR); err != nil {
		return err
	}

	logger := stageHarness.Logger

	fileName := randomFileName()
	fileContent := randomFileContent()

	logger.Debugf("Creating file %s in %s", fileName, DATA_DIR)
	logger.Debugf("File Content:\n%s", fileContent)
	err := createFileWith(DATA_DIR+fileName, fileContent)
	defer os.Remove(DATA_DIR + fileName)

	if err != nil {
		return err
	}

	err = testGetFileResponse(logger, fileName, fileContent)
	if err != nil {
		return err
	}

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

	resp, err := sendRequest(httpClient, req, logger)
	if err != nil {
		return err
	}

	if resp.StatusCode != 404 {
		return fmt.Errorf("Expected status code 404, got %d", resp.StatusCode)
	}

	return nil
}
