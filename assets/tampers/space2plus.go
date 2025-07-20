package main

import "strings"

type T struct{}

func (t T) Name() string {
	return "space2plus"
}

func (t T) Exec(payload string) string {
	return strings.ReplaceAll(payload, " ", "+")
}

func (t T) Desc() string {
	return "Replaces all spaces in the input string with pluses."
}

var Tamper T
