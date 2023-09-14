package request

import (
	"net/http"
	"strings"

	"github.com/Brum3ns/firefly/pkg/random"
)

type Insert struct {
	keyword string
	payload string
}

// Take the "Insert" structure and the "Request" structure
// Return the "Request" structure with the insert points replaced by the current payload.
func NewInsert(insert *Insert, req *RequestProperties) *RequestProperties {
	req.URL = insert.setURL(req.URL)
	req.Method = insert.setMethod(req.Method)
	req.PostBody = insert.setPostBody(req.PostBody)
	req.Headers = *insert.setHeaders(req.HeadersOriginal)
	return req
}

// Insert the payload based on the insert point (default=FUZZ) from user options to a string
func (ist *Insert) AddKeyword(source string) string {
	keyword := random.RandomInsert(ist.keyword)
	return strings.ReplaceAll(source, keyword, ist.payload)
}

func (ist Insert) setHeaders(sliceArry [][2]string) *http.Header {
	var headers = &http.Header{}
	for _, h := range sliceArry {
		hName := random.RandomInsert(h[0])
		hValue := random.RandomInsert(h[1])
		headers.Add(ist.AddKeyword(hName), ist.AddKeyword(hValue))
	}
	return headers
}

func (ist Insert) setURL(s string) string {
	keyword := random.RandomInsert(ist.keyword)
	return strings.ReplaceAll(s, keyword, URLNormalize(ist.payload))
}

func (ist Insert) setPostBody(s string) string {
	return ist.AddKeyword(s)
}

func (ist Insert) setMethod(s string) string {
	return ist.AddKeyword(s)
}
