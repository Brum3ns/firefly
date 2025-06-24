package rhttp

import (
	"fmt"
	"net/http"

	"github.com/projectdiscovery/rawhttp"
)

type Client struct {
	*rawhttp.Client
}

// New returns a wrapped client
func NewClient(options *rawhttp.Options) *Client {
	return &Client{
		Client: rawhttp.NewClient(options),
	}
}

func (client *Client) SendRawRequest(url string, rawRequest []byte) (*http.Response, error) {
	opts := *client.Options
	opts.CustomRawBytes = rawRequest
	resp, err := client.DoRawWithOptions(
		"",
		url,
		"",
		nil,
		nil,
		&opts,
	)
	if err != nil {
		return resp, fmt.Errorf("request failed, error: %v", err)
	}
	return resp, nil
}
