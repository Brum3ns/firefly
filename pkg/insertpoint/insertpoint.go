package insertpoint

import (
	"net/http"
	"strings"

	"github.com/Brum3ns/firefly/pkg/random"
)

type Insert struct {
	Keyword string
	Payload string
}

// Take the "Insert" structure and the "Request" structure
// Return the "Request" structure with the insert points replaced by the current payload.
func NewInsert(keyword, payload string) Insert {
	return Insert{
		Keyword: keyword,
		Payload: payload,
	}
}

// Insert the payload based on the insert point (default=FUZZ) from user options to a string
func (ist Insert) addKeyword(s string) string {
	return strings.ReplaceAll(random.RandomInsert(s), ist.Keyword, ist.Payload)
}

func (ist Insert) SetHeaders(sliceArry [][2]string) http.Header {
	var headers = http.Header{}
	for _, h := range sliceArry {
		hName := random.RandomInsert(h[0])
		hValue := random.RandomInsert(h[1])
		headers.Add(ist.addKeyword(hName), ist.addKeyword(hValue))
	}
	return headers
}

func (ist Insert) SetURL(s string) string {
	return strings.ReplaceAll(random.RandomInsert(s), ist.Keyword, normalizeURLstring(ist.Payload))
}

func (ist Insert) SetPostBody(s string) string {
	return ist.addKeyword(s)
}

func (ist Insert) SetMethod(s string) string {
	return ist.addKeyword(s)
}

// Normalize common characters in the URL into URL-encode:
func normalizeURLstring(s string) string {
	var (
		l_find        = []string{" ", "\t", "\n", "#", "&", "?"}
		l_URLEncodeTo = []string{"%20", "%09", "%0a", "%23", "%26", "%3F"}
	)
	for i := 0; i < len(l_URLEncodeTo); i++ {
		if strings.Contains(s, l_find[i]) {
			s = strings.ReplaceAll(s, l_find[i], l_URLEncodeTo[i])
		}
	}
	return s
}
