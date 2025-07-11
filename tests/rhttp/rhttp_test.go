package tests

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Brum3ns/firefly/pkg/rhttp"
	"github.com/projectdiscovery/rawhttp"
)

func Test_rhttp_test(t *testing.T) {
	rHttp, err := rhttp.NewRequest(
		rhttp.Config{
			Methods: []string{"GET"},
			Headers: map[string][]string{
				"Host":       {"2.example.com"},
				"X-Test":     {"1337"},
				"User-Agent": {"firefly v2.0"},
			},
			Url:                 "http://example.com/",
			URIPath:             "/?x=1",
			Body:                "y=2",
			RespectHSTS:         false,
			HeaderPresetBrowser: "firefox",
		},
		&rawhttp.Options{
			Timeout:                30 * time.Second,
			FollowRedirects:        true,
			MaxRedirects:           10,
			AutomaticHostHeader:    true,
			AutomaticContentLength: true,
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err := rHttp.MakeCoreRequests(); err != nil {
		log.Fatalln(err)
	}

	for _, rawReq := range rHttp.GetCoreRawRequests() {
		fmt.Println(string(rawReq.TemplateRawRequest))
	}
}
