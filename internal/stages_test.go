package internal

import (
	"os"
	"regexp"
	"testing"

	tester_utils_testing "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	os.Setenv("CODECRAFTERS_RANDOM_SEED", "1234567890")

	falseVar := false

	testCases := map[string]tester_utils_testing.TesterOutputTestCase{
		"init_failure": {
			UntilStageSlug:      "connect-to-port",
			CodePath:            "./test_helpers/scenarios/init/failure",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_timeout": {
			UntilStageSlug:      "connect-to-port",
			CodePath:            "./test_helpers/scenarios/init/timeout",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/timeout",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_success_cheat": {
			UntilStageSlug:      "connect-to-port",
			CodePath:            "./test_helpers/scenarios/init/success_cheat",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/success_cheat",
			NormalizeOutputFunc: normalizeTesterOutput,
			SkipAntiCheat:       &falseVar,
		},
		"init_success": {
			UntilStageSlug:      "connect-to-port",
			CodePath:            "./test_helpers/scenarios/init/success",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/init/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"pass_all": {
			UntilStageSlug:      "post-file",
			CodePath:            "./test_helpers/scenarios/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/pass_all",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	re, _ := regexp.Compile(`(\d{2}\/[A-Za-z]{3}\/\d{4} \d{2}:\d{2}:\d{2})`)
	return re.ReplaceAll(testerOutput, []byte("xx/xxx/xxxx xx:xx:xx"))
}
