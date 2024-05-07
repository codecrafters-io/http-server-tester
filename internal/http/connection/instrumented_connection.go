package http_connection

import (
	"net/http"

	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func defaultCallbacks(stageHarness *test_case_harness.TestCaseHarness, logPrefix string) HttpConnectionCallbacks {
	return HttpConnectionCallbacks{
		BeforeSendRequest: func(request *http.Request) {
			stageHarness.Logger.Infof("%s$ %s", logPrefix, HttpRequestToCurlString(request))
		},
		BeforeSendBytes: func(bytes []byte) {
			stageHarness.Logger.Debugf("%sSent bytes: %q", logPrefix, string(bytes))
		},
		AfterBytesReceived: func(bytes []byte) {
			stageHarness.Logger.Debugf("%sReceived bytes: %q", logPrefix, string(bytes))
		},
		AfterReadResponse: func(response http_parser.HTTPResponse) {
			stageHarness.Logger.Debugf("%sReceived response: %v", logPrefix, response.FormattedString())
		},
	}
}

func NewInstrumentedHttpConnection(stageHarness *test_case_harness.TestCaseHarness, addr string, clientIdentifier string) (*HttpConnection, error) {
	logPrefix := ""
	if clientIdentifier != "" {
		logPrefix = clientIdentifier + ": "
	}
	return NewHttpConnection(
		addr, defaultCallbacks(stageHarness, logPrefix),
	)
}
