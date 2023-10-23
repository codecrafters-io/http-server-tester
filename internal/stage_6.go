package internal

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	testerutils "github.com/codecrafters-io/tester-utils"
	logger "github.com/codecrafters-io/tester-utils/logger"
)

func testHandlesMultipleConcurrentConnections(stageHarness *testerutils.StageHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	randomInt := rand.Intn(2) + 2

	logger.Infof("Creating %d parallel connections", randomInt)
	conns := make([]net.Conn, randomInt)

	for i := 0; i < randomInt; i++ {
		logger.Debugf("Creating connection %d", i+1)
		conn, err := createTcpConn(TCP_DEST)
		if err != nil {
			return err
		}
		conns[i] = conn
	}
	for i := randomInt - 1; i >= 0; i-- {
		err := sendRequestDirectlyOverTcp(logger, conns[i], i+1)
		if err != nil {
			logFriendlyError(logger, err)
			return err
		}
	}
	for i := randomInt - 1; i >= 0; i-- {
		logger.Debugf("Closing connection %d", i+1)
		err := conns[i].Close()
		if err != nil {
			logFriendlyError(logger, err)
			return err
		}
	}

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

func sendRequestDirectlyOverTcp(logger *logger.Logger, conn net.Conn, i int) error {
	req := "GET / HTTP/1.1\r\n" + "\r\n\r\n"
	logger.Infof("Sending request on %d (status line): %q", i, getFirstLine(string(req)))
	logPrefix := ">>>"
	logger.Debugf("Sending request on %d: (Messages with %q prefix are part of this log)", i, logPrefix)
	logFriendlyHTTPMessage(logger, string(req), logPrefix)

	_, err := conn.Write([]byte(req))
	if err != nil {
		return err
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		return err
	}
	err = dumpResponse(logger, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != resp.StatusCode {
		return fmt.Errorf("Expected status code %d, got %d on connection %d", 200, resp.StatusCode, i)
	}
	defer resp.Body.Close()
	return nil
}
