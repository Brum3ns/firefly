package technique

import (
	"net/http"

	G "github.com/Brum3ns/FireFly/pkg/functions/globalVariables"
)

func General(body []byte, headers http.Header) (bool, []string) {
	var (
		reMatch     bool
		lst_headers = []string{}
	)

	//Check for mathing headers:
	if len(G.Lst_CheckHeaders) > 0 {
		for h := range headers {
			for _, H := range G.Lst_CheckHeaders {

				if h == H {
					lst_headers = append(lst_headers, h)
				}
			}
		}
	}

	return reMatch, lst_headers
}
