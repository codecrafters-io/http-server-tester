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
	ExecutableFileName:       "your_program.sh",
	LegacyExecutableFileName: "your_server.sh",
	TestCases: []tester_definition.TestCase{
		{
			Slug:     "at4",
			TestFunc: testConnects,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ia4",
			TestFunc: test200OK,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ih0",
			TestFunc: test404NotFound,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "cn2",
			TestFunc: testRespondWithContent,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "fs3",
			TestFunc: testRespondWithUserAgent,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ej5",
			TestFunc: testHandlesMultipleConcurrentConnections,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ap6",
			TestFunc: testGetFile,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "qv8",
			TestFunc: testPostFile,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "df4",
			TestFunc: testRespondWithContentEncoding,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ij8",
			TestFunc: testRespondWithCorrectContentEncoding,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "cr8",
			TestFunc: testRespondWithEncodedData,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "pe1",
			TestFunc: testPersistence1,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "pe2",
			TestFunc: testPersistence2,
			Timeout:  15 * time.Second,
		},
	},
}
