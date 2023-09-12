package internal

import (
	"math/rand"
	"strings"
)

var randomWords = []string{
	"humpty",
	"dumpty",
	"Horsey",
	"donkey",
	"yikes",
	"monkey",
	"Coo",
	"scooby",
	"dooby",
	"vanilla",
	"237",
	"Monkey",
}

func randomWord() string {
	return randomWords[rand.Intn(len(randomWords))]
}

func randomString(n int, joiner string) string {
	b := make([]string, n)

	for i := range b {
		b[i] = randomWord()
	}

	return strings.Join(b, joiner)
}

func randomUrl() string {
	size := rand.Intn(2) + 1
	return randomString(size, "-")
}

func randomUrlPath() string {
	return randomUrl()
}
