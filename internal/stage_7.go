package internal

import (
	"fmt"
	"io"
	"strings"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testFileListing(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	httpClient := NewHTTPClient()

	logger.Infof("Running stage 7")

	response, err := httpClient.Get("http://localhost:4221/files")
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

	expected := []string{
		"file1.txt",
		"file2.txt",
	}

	if !hasCommonElement(lines, expected) {
		return fmt.Errorf("Expected lines: %v, got: %v", expected, lines)
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
