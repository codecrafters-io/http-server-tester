package internal

import (
	executable "github.com/codecrafters-io/tester-utils/executable"
	logger "github.com/codecrafters-io/tester-utils/logger"
	testerutils "github.com/codecrafters-io/tester-utils"
)

type HTTPServerBinary struct {
	executable *executable.Executable
	logger     *logger.Logger
}

func NewHTTPServerBinary(stageHarness *testerutils.StageHarness) *HTTPServerBinary {
	b := &HTTPServerBinary{
		executable: stageHarness.Executable,
		logger:     stageHarness.Logger,
	}

	stageHarness.RegisterTeardownFunc(func() { b.Kill() })

	return b
}

func (b *HTTPServerBinary) Run(args ...string) error {
	b.logger.Debugf("Running program")
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
