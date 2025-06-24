package tests

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
)

func Foo(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
	return "works!", nil
}

func Test_faker_func(t *testing.T) {
	//gofakeit.Seed(0)

	// Register custom placeholder
	gofakeit.AddFuncLookup("fuzz", gofakeit.Info{
		Category:    "custom",
		Description: "Random name in lowercase",
		Example:     "john doe",
		Output:      "string",
		Generate:    Foo,
	})

	// Use custom tag in template
	fmt.Println(gofakeit.Generate("Payload: {fuzz}"))
}
