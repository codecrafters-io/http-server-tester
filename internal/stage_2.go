package internal

import (
	http_assertions "github.com/codecrafters-io/http-server-tester/internal/http/assertions"
	"github.com/codecrafters-io/http-server-tester/internal/http/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func test200OK(stageHarness *test_case_harness.TestCaseHarness) error {
	b := NewHTTPServerBinary(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	requestResponsePair, err := getBaseURLGetRequestResponsePair()
	if err != nil {
		return err
	}

	test_case := test_cases.SendRequestTestCase{
		Request:                   requestResponsePair.Request,
		Assertion:                 http_assertions.NewHTTPResponseAssertion(*requestResponsePair.Response),
		ShouldSkipUnreadDataCheck: false,
	}

	return test_case.Run(stageHarness, TCP_DEST, logger)
}
