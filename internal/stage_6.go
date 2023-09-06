package internal

import (
	"math/rand"
	"net"
	"time"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testHandlesMultipleConcurrentConnections(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	randomInt := rand.Intn(5) + 5

	logger.Infof("Creating %d parallel connections", randomInt)

	for i := 0; i < randomInt; i++ {
		conn, err := createTcpConn(TCP_DEST)
		if err != nil {
			return err
		}
		defer conn.Close()
	}

	httpClient := NewHTTPClient()
	requestWithStatus(httpClient, URL, 200, logger)

	return nil
}

func createTcpConn(destination string) (net.Conn, error) {
	retries := 0
	var conn net.Conn
	var err error

	for {
		conn, err = net.Dial("tcp", destination)
		if err != nil && retries > 15 {
			return nil, err
		} else if err == nil {
			return conn, nil
		} else {
			retries += 1
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
