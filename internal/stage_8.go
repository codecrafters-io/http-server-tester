package internal

import (
	"os"

	testerutils "github.com/codecrafters-io/tester-utils"
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

	return nil
}
