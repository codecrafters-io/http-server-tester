package internal

import (
	"time"

	testerutils "github.com/codecrafters-io/tester-utils"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatTestCases: []testerutils.TestCase{
		{
			Slug:                    "anti-cheat-1",
			TestFunc:                antiCheatBasic,
			Timeout:                 15 * time.Second,
		},
	},
	ExecutableFileName: "your_server.sh",
	TestCases: []testerutils.TestCase{
		{
			Slug:                    "connect-to-port",
			TestFunc:                testConnects,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "respond-with-200",
			TestFunc:                test200OK,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "respond-with-404",
			TestFunc:                test404NotFound,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "respond-with-content",
			TestFunc:                testRespondWithContent,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "parse-headers",
			TestFunc:                testRespondWithUserAgent,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "concurrent-connections",
			TestFunc:                testHandlesMultipleConcurrentConnections,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "get-file",
			TestFunc:                testGetFile,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "post-file",
			TestFunc:                testPostFile,
			Timeout:                 15 * time.Second,
		},
	},
}
