package internal

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func createFileWith(location string, content string) error {
	dirPath := filepath.Dir(location)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(location)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func testGetFileResponse(logger *testerutils.Logger, fileName string, fileContent string) error {
	httpClient := NewHTTPClient()

	req, err := http.NewRequest("GET", URL+"files/"+fileName, nil)
	if err != nil {
		return err
	}

	resp, err := sendRequest(httpClient, req, logger)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	if err := validateContent(*resp, fileContent); err != nil {
		return err
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
