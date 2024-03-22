package internal

import (
	"strings"

	executable "github.com/codecrafters-io/tester-utils/executable"
	logger "github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type HTTPServerBinary struct {
	executable *executable.Executable
	logger     *logger.Logger
}

func NewHTTPServerBinary(stageHarness *test_case_harness.TestCaseHarness) *HTTPServerBinary {
	b := &HTTPServerBinary{
		executable: stageHarness.Executable,
		logger:     stageHarness.Logger,
	}

	stageHarness.RegisterTeardownFunc(func() { b.Kill() })

	return b
}

func (b *HTTPServerBinary) Run(args ...string) error {
	b.logger.Debugf("Running program")
    if args == nil || len(args) == 0 {
        b.logger.Infof("$ ./your_server.sh")
    } else {
        b.logger.Infof("$ ./your_server.sh %s", strings.Join(args, " "))
    }
	if err := b.executable.Start(args...); err != nil {
		return err
	}

	return nil
}

func (b *HTTPServerBinary) HasExited() bool {
	return b.executable.HasExited()
}

func (b *HTTPServerBinary) Kill() error {
	b.logger.Debugf("Terminating program")
	if err := b.executable.Kill(); err != nil {
		b.logger.Debugf("Error terminating program: '%v'", err)
		return err
	}

	b.logger.Debugf("Program terminated successfully")
	return nil // When does this happen?
}
