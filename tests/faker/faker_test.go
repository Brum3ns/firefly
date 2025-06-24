package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/Brum3ns/firefly/pkg/faker"
)

func Test_gofakeit(t *testing.T) {
	const dataVal = `
GET /some/path?query={letter:0,10}&limit={number:1,1337}&cb={regex:\w{8}}#fragment
Host: http://example.com:80
Accept: */*
X-Ignore: ? # *

fatget[]=1337&email={email}&name={firstname}&nope={noexiting}`

	v, err := faker.Generate(dataVal)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(v)
}
