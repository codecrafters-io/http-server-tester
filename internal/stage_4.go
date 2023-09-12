package internal

import (
	"fmt"
	"io"
	"net/http"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testRespondWithContent(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	httpClient := NewHTTPClient()

	content := randomUrlPath()
	url := URL + "echo/" + content

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := sendRequest(httpClient, req, logger)

	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	if err := validateContent(*resp, content); err != nil {
		return err
	}

	return nil
}

func validateContent(resp http.Response, expected string) error {
	contentLength := len(expected)

	// test content-type

	receivedContentType := resp.Header.Get("Content-Type")
	if receivedContentType == "" {
		return fmt.Errorf("Content-Type header not present")
	}

	if receivedContentType != "text/plain" {
		return fmt.Errorf("Expected content type text/plain, got %s", receivedContentType)
	}

	// test content-length

	receivedContentLength := resp.Header.Get("Content-Length")
	if receivedContentLength == "" {
		return fmt.Errorf("Content-Length header not present")
	}

	if receivedContentLength != fmt.Sprintf("%d", contentLength) {
		return fmt.Errorf("Expected content length %d, got %s", contentLength, receivedContentLength)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != expected {
		return fmt.Errorf("Expected the content to be %s got %s", expected, body)
	}

	return nil
}
