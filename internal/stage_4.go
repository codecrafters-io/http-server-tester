package internal

import (
	"fmt"
	"io"
	"math/rand"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testRespondWithContent(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	httpClient := NewHTTPClient()

	contentLength := rand.Intn(5) + 5
	content := randSeq(contentLength)
	url := URL + "echo/" + content

	logger.Infof("Calling %s", url)

	response, err := httpClient.Get(url)
	if err != nil {
		logFriendlyError(logger, err)
		return fmt.Errorf("Failed to connect to server, err: '%v'", err)
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Expected status code 200, got %d", response.StatusCode)
	}

	// test content length header is contentLength long
	if cLen := response.Header.Get("Content-Length"); cLen != fmt.Sprintf("%d", contentLength) {
		return fmt.Errorf("Expected content length %d, got %s", contentLength, cLen)
	}

	// test body returns the passed path
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if string(body) != content {
		return fmt.Errorf("Expected the content to be %s got %s", content, body)
	}

	return nil
}
