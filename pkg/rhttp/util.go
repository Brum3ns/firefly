package rhttp

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/projectdiscovery/rawhttp"
	"github.com/projectdiscovery/rawhttp/client"
	urlutil "github.com/projectdiscovery/utils/url"
)

// DumpRequestRaw to string
func DumpRequestRaw(method, url, uripath string, headers map[string][]string, body io.Reader, options *rawhttp.Options) ([]byte, error) {
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
	if !hasHostHeader {
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

	req := toRequest(method, path, nil, headers, body, options)
	b := strings.Builder{}

	q := strings.Join(req.Query, "&")
	if len(q) > 0 {
		q = "?" + q
	}

	b.WriteString(fmt.Sprintf("%s %s%s %s"+client.NewLine, req.Method, req.Path, q, req.Version.String()))

	for _, header := range req.Headers {
		if header.Value != "" {
			b.WriteString(fmt.Sprintf("%s: %s"+client.NewLine, header.Key, header.Value))
		} else {
			b.WriteString(fmt.Sprintf("%s"+client.NewLine, header.Key))
		}
	}

	l := req.ContentLength()
	if req.AutomaticContentLength && l >= 0 {
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

func toRequest(method string, path string, query []string, headers map[string][]string, body io.Reader, options *rawhttp.Options) *client.Request {
	if len(options.CustomRawBytes) > 0 {
		return &client.Request{RawBytes: options.CustomRawBytes}
	}
	reqHeaders := toHeaders(headers)
	if len(options.CustomHeaders) > 0 {
		reqHeaders = options.CustomHeaders
	}

	return &client.Request{
		Method:  method,
		Path:    path,
		Query:   query,
		Version: client.HTTP_1_1,
		Headers: reqHeaders,
		Body:    body,
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
