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
			Slug:                    "init",
			Number:                  1,
			Title:                   "Can connect to a TCP server",
			TestFunc:                testConnects,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-1",
			Number:                  1,
			Title:                   "Respond with 200",
			TestFunc:                test200OK,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-2",
			Number:                  2,
			Title:                   "Respond with 404",
			TestFunc:                test404NotFound,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-3",
			Number:                  3,
			Title:                   "Handle multiple concurrent connections",
			TestFunc:                testHandlesMultipleConcurrentConnections,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
		{
			Slug:                    "stage-4",
			Number:                  4,
			Title:                   "List files",
			TestFunc:                testFileListing,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
	},
}
