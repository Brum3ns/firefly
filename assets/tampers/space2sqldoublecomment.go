package main

import "strings"

type T struct{}

func (t T) Name() string {
	return "space2sqldoublecomment"
}

func (t T) Exec(payload string) string {
	return strings.ReplaceAll(payload, " ", "/*/**/*/")
}

func (t T) Desc() string {
	return "Replaces all spaces in the payload with double SQL multi comments: '/*/**/*/'."
}

var Tamper T
