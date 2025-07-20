package main

type T struct{}

func (t T) Name() string {
	return "reverse"
}

func (t T) Exec(payload string) string {
	runes := []rune(payload)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (t T) Desc() string {
	return "Reverses the input string. Example: 'payload' -> 'daolyap'"
}

var Tamper T
