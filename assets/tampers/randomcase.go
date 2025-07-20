package main

import (
	"math/rand"
	"strings"
)

var (
	name = "randomcase"
)

type T struct{}

func (t T) Name() string {
	return name
}

func (t T) Exec(payload string) string {
	var result strings.Builder
	for _, char := range payload {
		s := string(char)
		if rand.Intn(2) == 0 {
			result.WriteString(strings.ToUpper(s))
		} else {
			result.WriteString(strings.ToLower(s))
		}
	}
	return result.String()
}

var Tamper T
