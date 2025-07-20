package main

import "strings"

var (
	name   = "uppercase"
	Tamper T
)

type T struct{}

func (t T) Name() string {
	return name
}

func (t T) Exec(payload string) string {
	return strings.ToUpper(payload)
}

func (t T) Desc() string {
	return "Converts the input string to uppercase. Example: input 'hello' -> output 'HELLO'."
}
