package internal

import (
	"time"

	testerutils "github.com/codecrafters-io/tester-utils"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatStages: []testerutils.Stage{
		{
			Slug:                    "anti-cheat-1",
			Title:                   "Anti-cheat 1",
			TestFunc:                antiCheatTest,
			ShouldRunPreviousStages: true,
		},
	},
	ExecutableFileName: "your_server.sh",
	Stages: []testerutils.Stage{
		{
			Slug:                    "connect-to-port",
			Number:                  1,
			Title:                   "Bind to a port",
			TestFunc:                testConnects,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "respond-with-200",
			Number:                  2,
			Title:                   "Respond with 200",
			TestFunc:                test200OK,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "respond-with-404",
			Number:                  3,
			Title:                   "Respond with 404",
			TestFunc:                test404NotFound,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "respond-with-content",
			Number:                  4,
			Title:                   "Respond with content",
			TestFunc:                testRespondWithContent,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "parse-headers",
			Number:                  5,
			Title:                   "Parse headers",
			TestFunc:                testRespondWithUserAgent,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "concurrent-connections",
			Number:                  6,
			Title:                   "Concurrent connections",
			TestFunc:                testHandlesMultipleConcurrentConnections,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "list-files",
			Number:                  7,
			Title:                   "List files",
			TestFunc:                testFileListing,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "get-file",
			Number:                  8,
			Title:                   "Get a file",
			TestFunc:                testGetFile,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "post-file",
			Number:                  9,
			Title:                   "Post a file",
			TestFunc:                testPostFile,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
	},
}
