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
			UntilStageSlug:      "at4",
			CodePath:            "./test_helpers/scenarios/init/failure",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_timeout": {
			UntilStageSlug:      "at4",
			CodePath:            "./test_helpers/scenarios/init/timeout",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/timeout",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_success_cheat": {
			UntilStageSlug:      "at4",
			CodePath:            "./test_helpers/scenarios/init/success_cheat",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/success_cheat",
			NormalizeOutputFunc: normalizeTesterOutput,
			SkipAntiCheat:       &falseVar,
		},
		"init_success": {
			UntilStageSlug:      "at4",
			CodePath:            "./test_helpers/scenarios/init/success",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/init/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"pass_all": {
			UntilStageSlug:      "qv8",
			CodePath:            "./test_helpers/scenarios/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/base/pass_all",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"compression_pass_all": {
			StageSlugs:          []string{"df4", "ij8", "cr8"},
			CodePath:            "./test_helpers/scenarios/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/compression/pass_all",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"persistence_pass_all": {
			StageSlugs:          []string{"ag9", "ul1", "kh7"},
			CodePath:            "./test_helpers/scenarios/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/persistence/pass_all",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"persistence_fail": {
			StageSlugs:          []string{"ag9"},
			CodePath:            "./test_helpers/scenarios/pass_base",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/persistence/fail",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	replacements := map[string][]*regexp.Regexp{
		"xx/xxx/xxxx xx:xx:xx": {regexp.MustCompile(`(\d{2}\/[A-Za-z]{3}\/\d{4} \d{2}:\d{2}:\d{2})`)},
		"gzip_encoded_data_11": {regexp.MustCompile(`\[stage-11\] .*Received bytes: .*`)},
		"gzip_encoded_data_10": {regexp.MustCompile(`\[stage-10\] .*Received bytes: .*`)},
		"gzip_encoded_data_9":  {regexp.MustCompile(`\[stage-9\] .*Received bytes: .*`)},
	}

	for replacement, regexes := range replacements {
		for _, regex := range regexes {
			testerOutput = regex.ReplaceAll(testerOutput, []byte(replacement))
		}
	}

	return testerOutput

}
