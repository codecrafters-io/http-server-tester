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
			Title:                   "Respond with 200",
			TestFunc:                test200OK,
			ShouldRunPreviousStages: true,
			Timeout:                 15 * time.Second,
		},
	},
}
