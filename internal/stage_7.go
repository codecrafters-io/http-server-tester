package internal

import (
	"fmt"
	"io"
	"os"
	"strings"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testFileListing(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	fileNames := []string{
		randSeq(FILENAME_SIZE),
		randSeq(FILENAME_SIZE),
	}
	err := createFiles(DATA_DIR, fileNames)
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	httpClient := NewHTTPClient()

	response, err := httpClient.Get(URL + "files")
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

	lines := strings.Split(string(body), "\n")

	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	if !hasCommonElement(lines, fileNames) {
		return fmt.Errorf("Expected lines: %v, got: %v", fileNames, lines)
	}

	err = deleteFiles(DATA_DIR, fileNames)
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	return nil
}

func contains(arr []string, target string) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}

func hasCommonElement(list1, list2 []string) bool {
	for _, item1 := range list1 {
		if contains(list2, item1) {
			return true
		}
	}
	return false
}

func createFiles(location string, fileNames []string) error {
	err := os.MkdirAll(location, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %v", err)
	}

	for _, fileName := range fileNames {
		err := createFileWith(DATA_DIR+fileName, randSeq(5))
		if err != nil {
			return fmt.Errorf("Failed to create file: %v", err)
		}
	}
	return nil
}

func deleteFiles(location string, fileNames []string) error {
	for _, fileName := range fileNames {
		err := os.Remove(location + fileName)
		if err != nil {
			return fmt.Errorf("Failed to delete file: %v", err)
		}
	}
	return nil
}
