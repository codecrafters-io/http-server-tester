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
			Slug:                    "stage-1",
			Number:                  1,
			Title:                   "Can connect to a TCP server",
			TestFunc:                testConnects,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-2",
			Number:                  2,
			Title:                   "Respond with 200",
			TestFunc:                test200OK,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-3",
			Number:                  3,
			Title:                   "Respond with 404",
			TestFunc:                test404NotFound,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-4",
			Number:                  4,
			Title:                   "Respond with content",
			TestFunc:                testRespondWithContent,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-5",
			Number:                  5,
			Title:                   "Header Parsing",
			TestFunc:                testRespondWithUserAgent,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-6",
			Number:                  6,
			Title:                   "Handle multiple concurrent connections",
			TestFunc:                testHandlesMultipleConcurrentConnections,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-7",
			Number:                  7,
			Title:                   "List files",
			TestFunc:                testFileListing,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-8",
			Number:                  8,
			Title:                   "Get file",
			TestFunc:                testGetFile,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
	},
}
