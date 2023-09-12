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

	fileName := randSeq(FILENAME_SIZE)
	fileContent := randSeq(1000)

	err := postFile(fileName, fileContent)
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	err = testGetFileResponse(logger, fileName, fileContent)
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	err = validateFile(fileName, fileContent)
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	err = os.Remove(DATA_DIR + fileName)
	if err != nil {
		logFriendlyError(logger, err)
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

func validateFile(fileName string, fileContent string) error {
	onDiskContent, err := os.ReadFile(DATA_DIR + fileName)
	if err != nil {
		return fmt.Errorf("Error reading file: %v", err)
	}
	if fileContent != string(onDiskContent) {
		return fmt.Errorf("Expected the content to be %s got %s", fileContent, onDiskContent)
	}
	return nil
}
