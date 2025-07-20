package main

import "strings"

type T struct{}

func (t T) Name() string {
	return "space2urlencode"
}

func (t T) Exec(payload string) string {
	return strings.ReplaceAll(payload, " ", "%20")
}

func (t T) Desc() string {
	return "Replaces spaces in the input string with '%20'. Example: 'Hello World' -> 'Hello%20World'."
}

var Tamper T
