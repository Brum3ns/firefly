package tests

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/Brum3ns/firefly/pkg/httpnode"
)

func Test_httpnode_detect(t *testing.T) {
	var respBodies = []string{
		`{"key": "value"}`,
		`dsa<html><body>Hello, World!</body></html>`,
		`<note><to>Tove</to><from>Jani</from><heading>Reminder</heading><body>Don't forget me this weekend!</body></note>`,
	}

	var respHeaders = []http.Header{
		{"Content-Type": []string{"application/json"}},
		{"Content-Type": []string{"application/xml"}},
		{"Content-Type": []string{"text/html"}},
	}

	for _, body := range respBodies {
		log.Println(httpnode.DetectResponseType(body))
	}
	fmt.Println("---")
	for _, headers := range respHeaders {
		log.Println(httpnode.DetectContentType(headers, false))
	}

}
