package rhttp

import (
	"slices"
	"strings"

	"github.com/projectdiscovery/rawhttp/client"
)

var (
	HEADERS_SUPPORTED_BROWSERS = []string{"chrome", "firefox", "safari"}

	HEADERS_CHROME = []client.Header{
		{Key: "sec-ch-ua", Value: "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\""},
		{Key: "sec-ch-ua-mobile", Value: "?0"},
		{Key: "sec-ch-ua-platform", Value: "\"macOS\""},
		{Key: "Upgrade-Insecure-Requests", Value: "1"},
		{Key: "User-Agent", Value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"},
		{Key: "Accept", Value: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		{Key: "Sec-Fetch-Site", Value: "none"},
		{Key: "Sec-Fetch-Mode", Value: "navigate"},
		{Key: "Sec-Fetch-User", Value: "?1"},
		{Key: "Sec-Fetch-Dest", Value: "document"},
		{Key: "Accept-Encoding", Value: "gzip, deflate, br"},
		{Key: "Accept-Language", Value: "en-US,en;q=0.9"},
	}

	HEADERS_FIREFOX = []client.Header{
		{Key: "User-Agent", Value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 14.7; rv:128.0) Gecko/20100101 Firefox/128.0"},
		{Key: "Accept", Value: "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
		{Key: "Accept-Language", Value: "en-US,en;q=0.5"},
		{Key: "Accept-Encoding", Value: "gzip, deflate, br"},
		{Key: "Connection", Value: "keep-alive"},
		{Key: "Upgrade-Insecure-Requests", Value: "1"},
		{Key: "Sec-Fetch-Dest", Value: "document"},
		{Key: "Sec-Fetch-Mode", Value: "navigate"},
		{Key: "Sec-Fetch-Site", Value: "none"},
		{Key: "Sec-Fetch-User", Value: "?1"},
	}

	HEADERS_SAFARI = []client.Header{
		{Key: "User-Agent", Value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_7_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.4.1 Safari/605.1.15"},
		{Key: "Accept", Value: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		{Key: "Accept-Language", Value: "en-US,en;q=0.9"},
		{Key: "Accept-Encoding", Value: "gzip, deflate, br"},
		{Key: "Connection", Value: "keep-alive"},
		{Key: "Upgrade-Insecure-Requests", Value: "1"},
	}
)

func getHeaderPreset(browser string) []client.Header {
	switch strings.ToLower(browser) {
	case "chrome":
		return HEADERS_CHROME
	case "firefox":
		return HEADERS_FIREFOX
	case "safari":
		return HEADERS_SAFARI

	}
	return []client.Header{}
}

func validHeaderPreset(browser string) bool {
	return slices.Contains(HEADERS_SUPPORTED_BROWSERS, strings.ToLower(browser))
}
