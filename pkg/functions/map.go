package functions

import (
	"fmt"
)

func MapAmount(mp map[string]string) map[string]int {
	m := make(map[string]int)
	for k, _ := range mp {
		m[fmt.Sprintf("%v", k)] += 1
	}
	return m
}
