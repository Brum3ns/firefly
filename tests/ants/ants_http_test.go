package main

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/Brum3ns/firefly/pkg/rhttp"
	"github.com/google/uuid"
	"github.com/panjf2000/ants/v2"
	"github.com/projectdiscovery/rawhttp"
)

const threads = 50

var client = rhttp.NewClient(&rawhttp.Options{
	Timeout:                7,
	FollowRedirects:        true,
	MaxRedirects:           12,
	AutomaticHostHeader:    false,
	AutomaticContentLength: false,
	ForceReadAllBody:       true,
})

var rawRequest = []byte(
	"GET / HTTP/1.1\r\n" +
		"Host: localhost:1337\r\n" +
		"User-Agent: Test ants x rawhttp\r\n" +
		"Connection: close\r\n" +
		"\r\n",
)

func sendRequest(i interface{}) {
	url := "http://localhost:1337/?cb=" + uuid.NewString()

	resp, err := client.SendRawRequest(url, rawRequest)

	if err != nil {
		log.Fatalln("Request error:", err)
	}
	fmt.Println(resp.StatusCode)

	resp.Body.Close()
}

func Test_ants_http(t *testing.T) {
	var wg sync.WaitGroup

	pool, _ := ants.NewPool(threads)
	defer pool.Release()

	requests := 1000
	for i := 0; i < requests; i++ {
		wg.Add(1)
		pool.Submit(func() {
			sendRequest(i)
			wg.Done()
		})
		//pool.Invoke(i)
	}

	wg.Wait()
}
