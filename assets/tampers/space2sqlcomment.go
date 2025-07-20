package main

import "strings"

type T struct{}

func (t T) Name() string {
	return "space2sqlcomment"
}

func (t T) Desc() string {
	return "Replaces all spaces in the input string with a SQL multiline comment."
}

func (t T) Exec(payload string) string {
	return strings.ReplaceAll(payload, " ", "/**/")
}

var Tamper T
