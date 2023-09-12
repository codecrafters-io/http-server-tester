package internal

import (
	"bytes"
	"fmt"
	"os"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testPostFile(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	fileName := randomFileName()
	fileContent := randomFileContent()

	err := postFile(fileName, fileContent)
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}
	defer os.Remove(DATA_DIR + fileName)

	err = testGetFileResponse(logger, fileName, fileContent)
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	err = validateFile(fileName, fileContent)
	if err != nil {
		return err
	}

	return nil
}

func postFile(fileName string, fileContent string) error {
	httpClient := NewHTTPClient()
	response, err := httpClient.Post(URL+"files/"+fileName, "text/plain", bytes.NewBufferString(fileContent))
	if err != nil {
		return fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}
	if response.StatusCode != 201 {
		return fmt.Errorf("Expected status code 201, got %d", response.StatusCode)
	}
	return nil
}
