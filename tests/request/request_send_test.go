package tests

import (
	"fmt"
	"log"
	"net/http/httputil"
	"strings"
	"testing"
	"time"

	"github.com/Brum3ns/firefly/pkg/rhttp"
	"github.com/projectdiscovery/rawhttp"
)

func Test_request_send(t *testing.T) {
	const url = "http://127.0.0.1:1337/"
	// Configure HTTP request
	req, err := rhttp.NewRequest(rhttp.Config{
		Methods: []string{"GET"},
		Headers: rhttp.MakeHeaders([]string{
			"X-Dup: x",
			"x-Dup: x",
			"x-dup: x",
			"X-Test: y",
		}),
		Url:         url,
		URIPath:     "/",
		Body:        "x=1",
		RespectHSTS: false,
	},
		&rawhttp.Options{
			Timeout:                time.Duration(3000) * time.Millisecond,
			FollowRedirects:        false,
			MaxRedirects:           0,
			AutomaticHostHeader:    false,
			AutomaticContentLength: false,
			//CustomHeaders:          opt.Headers,
			ForceReadAllBody: true,
			//CustomRawBytes:         "",
			Proxy:            "",
			ProxyDialTimeout: time.Duration(0) * time.Millisecond,
			//SNI:                    "",
			//FastDialer:             "",
		})
	if err != nil {
		log.Fatalln(err)
	}

	// Make raw request(s)
	if err := req.MakeCoreRawRequests(); err != nil {
		log.Fatal(err)
	}

	for _, rawReq := range req.GetCoreRawRequests() {
		resp, err := req.Client.SendRawRequest(url, rawReq)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(rawReq))

		fmt.Println(strings.Repeat("-", 33))

		b, err := httputil.DumpResponse(resp, false)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(b))
	}

}
