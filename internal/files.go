package internal

import (
	"fmt"
	"os"
	"path/filepath"

	logger "github.com/codecrafters-io/tester-utils/logger"
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

func validateFile(logger *logger.Logger, fileName string, fileContent string) error {
	logger.Debugf("Validating file `%s` exists on disk", fileName)
	onDiskContent, err := os.ReadFile(DATA_DIR + fileName)
	if err != nil {
		return fmt.Errorf("Error reading file: %v", err)
	}
	logger.Debugf("Validating file `%s` content", fileName)
	if fileContent != string(onDiskContent) {
		logger.Errorf("Expected file content: %q", fileContent)
		logger.Errorf("Received file content: %q", onDiskContent)
		return fmt.Errorf("File content does not match")
	}
	return nil
}
