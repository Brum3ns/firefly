package main

import (
	"regexp"
	"strings"
)

var (
	re = regexp.MustCompile(`\s+`)
)

type T struct{}

// Replace all whitespace sequences
func (t T) Name() string {
	return "whitespacecollapse"
}

func (t T) Exec(payload string) string {
	collapsed := re.ReplaceAllString(payload, " ")
	return strings.TrimSpace(collapsed)
}

func (t T) Desc() string {
	return "Collapses all whitespace sequences into a single space. Example: 'Hello   World' -> 'Hello World'"
}

var Tamper T
