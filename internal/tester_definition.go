package internal

import (
	"time"

	"github.com/codecrafters-io/tester-utils/tester_definition"
)

var testerDefinition = tester_definition.TesterDefinition{
	AntiCheatTestCases: []tester_definition.TestCase{
		{
			Slug:     "anti-cheat-1",
			TestFunc: antiCheatBasic,
			Timeout:  15 * time.Second,
		},
	},
	ExecutableFileName: "your_server.sh",
	TestCases: []tester_definition.TestCase{
		{
			Slug:     "connect-to-port",
			TestFunc: testConnects,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "respond-with-200",
			TestFunc: test200OK,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "respond-with-404",
			TestFunc: test404NotFound,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "respond-with-content",
			TestFunc: testRespondWithContent,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "parse-headers",
			TestFunc: testRespondWithUserAgent,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "concurrent-connections",
			TestFunc: testHandlesMultipleConcurrentConnections,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "get-file",
			TestFunc: testGetFile,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "post-file",
			TestFunc: testPostFile,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "compression-content-encoding",
			TestFunc: testRespondWithContentEncoding,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "compression-multiple-schemes",
			TestFunc: testRespondWithCorrectContentEncoding,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "compression-gzip",
			TestFunc: testRespondWithEncodedData,
			Timeout:  15 * time.Second,
		},
	},
}
