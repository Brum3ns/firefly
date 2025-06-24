package tests

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/projectdiscovery/rawhttp"
)

func Test_rawhttp(t *testing.T) {
	var (
		method  = "GET"
		url     = "http://127.0.0.1:1337/"
		headers = map[string][]string{
			//"x-test;": {"1"},
		}
		body         = strings.NewReader("x=1")
		rawHTTPBytes = []byte("GET / HTTP/1.1\r\nHost: 127.0.0.1:1337\r\nX-Mal;\r\n\r\n")
		//proxy = "http://127.0.0.1:8080/"
	)

	opt := &rawhttp.Options{
		Timeout:                5 * time.Second,
		FollowRedirects:        true,
		MaxRedirects:           3,
		AutomaticHostHeader:    false,
		AutomaticContentLength: true,
		CustomRawBytes:         rawHTTPBytes,
	}

	rawBytes, err := rawhttp.DumpRequestRaw(method,
		url,
		"",
		headers,
		body,
		opt,
	)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(rawBytes))
	fmt.Println("---")

	resp, err := rawhttp.DoRawWithOptions(
		method,
		url,
		"",
		headers,
		body,
		opt,
	)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp)
}
