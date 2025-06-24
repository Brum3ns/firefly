package tests

import (
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/projectdiscovery/rawhttp"
)

func Test_rawhttp_hts(t *testing.T) {
	var (
		method  = "POST"
		url     = "https://0a90007603ee4a9681c539b6006e00ad.web-security-academy.net/"
		headers = map[string][]string{
			//"x-test;": {"1"},
		}
		body         = strings.NewReader("x=1")
		rawHTTPBytes = []byte("POST / HTTP/1.1\r\nHost: 0a90007603ee4a9681c539b6006e00ad.web-security-academy.net\r\nConnection: keep-alive\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 6\r\nTransfer-Encoding: chunked\r\n\r\n0\r\n\r\nG\r\n\r\n")
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
	// Remember to close the body
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(string(bodyBytes))
}
