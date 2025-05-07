package internal

import (
	"strings"

	"github.com/codecrafters-io/tester-utils/logger"
)

func logFriendlyError(logger *logger.Logger, err error) {
	if strings.Contains(err.Error(), "EOF") {
		logger.Infof("Hint: EOF is short for 'end of file'. This usually means that your program either:")
		logger.Infof(" (a) didn't send a complete response, or")
		logger.Infof(" (b) closed the connection early")
	}

	if strings.Contains(err.Error(), "connection reset by peer") {
		logger.Infof("Hint: 'connection reset by peer' usually means that your program closed the connection before sending a complete response.")
	}

	if strings.Contains(err.Error(), "connection refused") {
		logger.Infof("Hint: 'connection refused' usually means that your server is not running.")
	}
}
