package tests

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/Brum3ns/firefly/pkg/httpnode"
)

func Test_httpprepare_htmlnode(t *testing.T) {
	const url = "https://example.com/"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	htmlnode := httpnode.GetHTMLNode(string(bodyResp))

	htmlJson, err := httpnode.HTMLNodeToJson(htmlnode)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(htmlJson))

}
