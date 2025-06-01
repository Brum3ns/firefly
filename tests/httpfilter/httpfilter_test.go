package tests

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Brum3ns/firefly/pkg/httpfilter"
)

func config() (httpfilter.Filter, error) {
	headersIgnore := http.Header{}

	filter, err := httpfilter.NewFilter(httpfilter.Config{
		Mode:                  "",             //"and",                                           // OK
		HeaderRegex:           "",             //"Da[Tt]e+",                                      // OK
		BodyRegex:             "",             //"i[a-z] mi[s]{2}ing",                            // OK
		StatusCodes:           []string{},     //[]string{" ++100", "200", "  100-300", "--201"}, // OK
		ResponseSizes:         []string{},     //[]string{"100-500"},                             // OK
		WordCounts:            []string{"47"}, //[]string{"47 "},                                 // OK
		LineCounts:            []string{},     //[]string{"--50", "100-500  "},                   // OK
		ResponseTimesMillisec: []string{},     //[]string{" ++0.0015   "},                        // OK
		Header:                headersIgnore,  // OK
	})

	return filter, err
}

func Test_HttpFilter(t *testing.T) {
	filter, err := config()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Do basic requests:
	for i := 0; i < 10; i++ {
		url := "http://localhost:1337/"

		timer := time.Now()

		resp, _ := http.Get(url)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Response body error", err)
		}
		respTimeSeconds := time.Since(timer).Seconds()
		bodyString := string(bodyBytes[:])
		bodySize := len(bodyBytes)
		wordCount := len(strings.Fields(bodyString))
		lineCount := len(strings.Split(bodyString, "\n"))

		filterRespone := httpfilter.Response{
			Body:         bodyBytes,
			StatusCode:   resp.StatusCode,
			ResponseSize: bodySize,
			LineCount:    lineCount,
			WordCount:    wordCount,
			ResponseTime: respTimeSeconds,
			Headers:      resp.Header,
		}

		// Filter requests
		fmt.Printf("[%d] %s - %d, Size:%d, LC:%d, WC:%d, Time:%f",
			i,
			resp.Request.URL,
			resp.StatusCode,
			bodySize,
			lineCount,
			wordCount,
			respTimeSeconds,
		)
		// fmt.Printf("|- Headers: %+v", resp.Header)
		print("\n")

		if filter.Run(filterRespone) {
			fmt.Println("[\033[1;36m>\033[0m] Filtered ID:", i)
		}
	}
}
