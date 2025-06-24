package tests

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
)

func Test_list_gofakeit(t *testing.T) {
	for name := range gofakeit.FuncLookups {
		fmt.Println(name)
	}
}
