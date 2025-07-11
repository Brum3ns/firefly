package rhttp

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/projectdiscovery/rawhttp"
	"github.com/projectdiscovery/rawhttp/client"
	urlutil "github.com/projectdiscovery/utils/url"
)

var (
	HTTP_1_1 = client.Version{
		Major: 1,
		Minor: 1,
	}
	HTTP_1_0 = client.Version{
		Major: 1,
		Minor: 0,
	}
	HTTP_0_9 = client.Version{
		Major: 0,
		Minor: 9,
	}
)

// DumpRequestRaw to string
func DumpRequestRaw(version client.Version, method, url, uripath string, headers map[string][]string, headersBrowserPreset string, body io.Reader, options *rawhttp.Options) ([]byte, error) {
	//headers := makeRawHeaders(headersRaw)
	if len(options.CustomRawBytes) > 0 {
		return options.CustomRawBytes, nil
	}
	if headers == nil {
		headers = make(map[string][]string)
	}
	u, err := urlutil.ParseURL(url, true)
	if err != nil {
		return nil, err
	}

	// Handle only if host header is missing
	_, hasHostHeader := headers["Host"]
	if options.AutomaticHostHeader && !hasHostHeader {
		host := u.Host
		headers["Host"] = []string{host}
	}

	// standard path
	path := u.Path
	if path == "" {
		path = "/"
	}
	if !u.Params.IsEmpty() {
		path += "?" + u.Params.Encode()
	}
	// override if custom one is specified
	if uripath != "" {
		path = uripath
	}

	req := toRequest(version, method, path, nil, headers, body, options)
	b := strings.Builder{}

	q := strings.Join(req.Query, "&")
	if len(q) > 0 {
		q = "?" + q
	}

	// Write HTTP first-line
	b.WriteString(fmt.Sprintf("%s %s%s %s"+client.NewLine, req.Method, req.Path, q, req.Version.String()))

	// Write HTTP headers
	for _, header := range makeHeaderOrder(req.Headers, headersBrowserPreset) {
		if header.Value != "" {
			b.WriteString(fmt.Sprintf("%s: %s"+client.NewLine, header.Key, header.Value))
		} else {
			b.WriteString(fmt.Sprintf("%s"+client.NewLine, header.Key))
		}
	}

	l := req.ContentLength()
	if options.AutomaticContentLength && l >= 0 {
		b.WriteString(fmt.Sprintf("Content-Length: %d"+client.NewLine, l))
	}

	b.WriteString(client.NewLine)

	if req.Body != nil {
		var buf bytes.Buffer
		tee := io.TeeReader(req.Body, &buf)
		body, err := io.ReadAll(tee)
		if err != nil {
			return nil, err
		}
		b.Write(body)
	}

	return []byte(b.String()), nil
}

func makeHeaderOrder(headers []client.Header, browser string) []client.Header {
	var headerOrder []client.Header

	if header, ok := GetHeader("Host", headers); ok {
		headerOrder = append(headerOrder, header)
		headers = RemoveHeader(header.Key, headers)
	}

	// Set the Host header first
	for _, headerPreset := range getHeaderPreset(browser) {
		if customHeader, ok := GetHeader(headerPreset.Key, headers); ok {
			headerOrder = append(headerOrder, customHeader)
			// Remove header from headers, since it was added as the preset headers
			headers = RemoveHeader(customHeader.Key, headers)

		} else {
			headerOrder = append(headerOrder, headerPreset)
		}
	}
	// Add the custom headers
	headerOrder = append(headerOrder, headers...)
	return headerOrder
}

// Check a string if it contains a HTTP version version
func MakeHTTPVersion(version string) client.Version {
	switch strings.ToLower(version) {
	case "1.0":
		return HTTP_1_0
	case "0.9":
		return HTTP_0_9
	default: /*1.1*/
		return HTTP_1_1
	}
}

// Take a list of raw headers and return a map where the items represent the amount of the header not value
// Note : to make it work with the rawhttp package's function: DumpRequestRaw
func MakeHeaders(headersRaw []string) map[string][]string {
	headers := make(map[string][]string)
	for _, header := range headersRaw {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			if values, ok := headers[parts[0]]; ok {
				values = append(values, parts[1])
				headers[parts[0]] = values
			} else {
				headers[parts[0]] = []string{parts[1]}
			}
		}
	}
	return headers
}

func toRequest(version client.Version, method string, path string, query []string, headers map[string][]string, body io.Reader, options *rawhttp.Options) *client.Request {
	if len(options.CustomRawBytes) > 0 {
		return &client.Request{RawBytes: options.CustomRawBytes}
	}
	reqHeaders := toHeaders(headers)
	if len(options.CustomHeaders) > 0 {
		reqHeaders = options.CustomHeaders
	}

	return &client.Request{
		Method:                 method,
		Path:                   path,
		Query:                  query,
		Version:                version,
		AutomaticContentLength: options.AutomaticContentLength,
		AutomaticHost:          options.AutomaticHostHeader,
		Headers:                reqHeaders,
		Body:                   body,
	}
}

const MaxResponseReadSizeDecompress = 10 * 1024 * 1024

func toHeaders(h map[string][]string) []client.Header {
	var r []client.Header
	for k, v := range h {
		for _, v := range v {
			r = append(r, client.Header{Key: k, Value: v})
		}
	}
	return r
}

func GetJitter(baseMs int) time.Duration {
	if baseMs <= 0 {
		return 0
	}
	// Set jitter range (e.g., 20% of base)
	jitterPercent := 0.2
	jitterRange := int(float64(baseMs) * jitterPercent)

	// Random jitter in the range [-jitterRange, +jitterRange]
	jitter := rand.Intn(2*jitterRange+1) - jitterRange

	totalDelay := baseMs + jitter
	return time.Duration(totalDelay) * time.Millisecond
}

// HeaderToString serializes an http.Header to its wire‑format string.
// Each key/value pair becomes one or more lines of the form:
//
//	Key: Value\r\n
//
// We sort the header names deterministically to make logging/tests repeatable.
// Note: RFC‑compliant line ending
func HeaderToString(header http.Header) string {
	var b strings.Builder

	// Sort the keys so output is stable
	keys := make([]string, 0, len(header))
	for k := range header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		for _, v := range header[k] {
			// RFC‑compliant line ending
			b.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
	}
	return b.String()
}

func GetHeader(headerKey string, headers []client.Header) (client.Header, bool) {
	for _, header := range headers {
		if header.Key == headerKey {
			return header, true
		}
	}
	return client.Header{}, false
}

func RemoveHeader(removeHeader string, headers []client.Header) []client.Header {
	var h []client.Header
	for _, header := range headers {
		if header.Key != removeHeader {
			h = append(h, header)
		}
	}
	return h
}
