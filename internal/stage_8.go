package internal

import (
	"fmt"
	"io"
	"os"
	"strings"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testGetFile(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	fileName := randSeq(FILENAME_SIZE)
	fileContent := randSeq(1000)
	err := createFileWith(DATA_DIR+fileName, fileContent)
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to create file: %v", err)
	}

	httpClient := NewHTTPClient()

	response, err := httpClient.Get(URL + "files" + "/" + fileName)
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Expected status code 200, got %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if strings.TrimSpace(string(body)) != fileContent {
		return fmt.Errorf("Expected the content to be %s got %s", fileContent, body)
	}

	err = os.Remove(DATA_DIR + fileName)
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	return nil
}
