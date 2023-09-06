package internal

import (
	"regexp"
	"testing"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func TestStages(t *testing.T) {
	testCases := map[string]tester_utils.TesterOutputTestCase{
		"init_failure": {
			StageName:           "connect-to-port",
			CodePath:            "./test_helpers/scenarios/init/failure",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_timeout": {
			StageName:           "connect-to-port",
			CodePath:            "./test_helpers/scenarios/init/timeout",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/timeout",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_success": {
			StageName:           "connect-to-port",
			CodePath:            "./test_helpers/scenarios/init/success",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/init/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	re, _ := regexp.Compile(`(\d{2}\/[A-Za-z]{3}\/\d{4} \d{2}:\d{2}:\d{2})`)
	return re.ReplaceAll(testerOutput, []byte("xx/xxx/xxxx xx:xx:xx"))
}
