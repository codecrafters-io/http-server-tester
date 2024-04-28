package internal

import (
	"github.com/codecrafters-io/tester-utils/random"
	"math/rand"
	"strings"
)

func randomString(n int, joiner string) string {
	b := make([]string, n)

	for i := range b {
		b[i] = random.RandomWord()
	}

	return strings.Join(b, joiner)
}

func randomAnything() string {
	size := rand.Intn(2) + 1
	return random.RandomWord() + "/" + randomString(size, "-")
}

func randomUrlPath() string {
	return random.RandomWord()
}

func randomUserAgent() string {
	return randomAnything()
}

func randomFileName() string {
	return randomString(4, "_")
}

func randomFileNameWithPrefix(prefix string) string {
	return prefix + randomFileName()
}

func randomFileContent() string {
	return randomString(8, " ")
}
