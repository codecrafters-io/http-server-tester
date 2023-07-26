package internal

import (
	"regexp"
	"testing"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func TestStages(t *testing.T) {
	testCases := map[string]tester_utils.TesterOutputTestCase{
		"init_failure": {
			StageName:           "init",
			CodePath:            "./test_helpers/scenarios/init/failure",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_timeout": {
			StageName:           "init",
			CodePath:            "./test_helpers/scenarios/init/timeout",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/timeout",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_success": {
			StageName:           "init",
			CodePath:            "./test_helpers/scenarios/init/success",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/init/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	re, _ := regexp.Compile("read tcp 127.0.0.1:\\d+->127.0.0.1:4221: read: connection reset by peer")
	return re.ReplaceAll(testerOutput, []byte("read tcp 127.0.0.1:xxxxx->127.0.0.1:4221: read: connection reset by peer"))
}
