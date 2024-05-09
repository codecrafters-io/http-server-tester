package internal

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

const URL = "http://127.0.0.1:4221/"
const TCP_DEST = "127.0.0.1:4221"
const DATA_DIR = "/tmp/data/codecrafters.io/http-server-tester/"
const FILENAME_SIZE = 40

func decodeGZIP(encodedString []byte) ([]byte, error) {
	reader := bytes.NewReader([]byte(encodedString))
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("Failed to decompress data: %v", err)
	}
	defer gzipReader.Close()

	decompressedData, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("Failed to decompress data: %v", err)
	}

	return decompressedData, nil
}
