package main

import "strings"

var (
	name = "lowercase"
)

type T struct{}

func (t T) Name() string {
	return name
}

func (t T) Exec(payload string) string {
	return strings.ToLower(payload)
}

var Tamper T
