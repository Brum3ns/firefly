package tests

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/Brum3ns/firefly/pkg/httpnode"
)

func Test_httpprepare_header(t *testing.T) {
	const url = "https://example.com/"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	headernode := httpnode.GetHeaderNode(resp.Header)

	headerJson, err := httpnode.HeaderNodeToJson(headernode)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(headerJson))

}
