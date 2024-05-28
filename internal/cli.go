package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
)

func RunCLI(env map[string]string) int {
	return testerutils.RunCLI(env, testerDefinition)
}
