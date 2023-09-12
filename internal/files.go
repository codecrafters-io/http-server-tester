package internal

import (
	"fmt"
	"net/http"
	"os"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func createFileWith(location string, content string) error {
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	if err := validateContent(*resp, fileContent); err != nil {
		return err
	}

	return nil
}
