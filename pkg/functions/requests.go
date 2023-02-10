package functions

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	G "github.com/Brum3ns/firefly/pkg/functions/globalVariables"
)

/*
func TotalRequests(u, m []string) int {
	Total := (len(u) * len(m) * G.Amount_Item)

	return Total
} */

/**Make the http.header (map) into a string. This is usefully for fast regex and diff checks later on*/
func HeadersToStr(headers http.Header) string {
	strHeaders := ""
	for k, v := range headers {
		strHeaders += fmt.Sprintf("%s: %s\n", k, strings.Join(v, " "))
	}

	return strHeaders
}

/*
func Hostname(u string) string {
	u = RegexBetween(u, "^[a-zA-Z]+://*(.*?)*($|/)")
	if u[len(u)-1:] == "/" {
		u = u[0 : len(u)-1]
	}

	return u
} */

//[TODO] - Improve and better storage
func RawData(data string) (map[string][]string, string) {
	var (
		method, url, path string
		RAW               = make(map[string][]string)
		v                 = true
	)

	lst_raw := strings.Split(data, "\n")

	for n, i := range lst_raw {
		s := strings.Split(i, " ")

		fmt.Println(n, i, s, "::", len(s)) //DEBUG

		//Extract "HTTP Method" and "Endpoint"
		if n == 0 {
			method = string(s[0])
			path = string(s[1])
			continue
		}

		//Extract Host header:
		host_Header := strings.ToLower(s[0])
		if host_Header == "host:" && v {
			url = string(s[1])
			v = false
			continue
		}

	}
	//Clear the first line Ex: "GET /endpoint HTTP/1.1"
	lst_raw = append(lst_raw[:0], lst_raw[1:]...)

	// If the secound last line is "\n" and last line isen't. Add the body to request:

	body := ""
	if lst_raw[len(lst_raw)-2] == "" && len(lst_raw[len(lst_raw)-1]) > 0 {
		fmt.Println("TRUE")
		body = lst_raw[len(lst_raw)-1]
	}

	//If not URL is detected there is not host header detected:
	if len(url) <= 0 {
		return nil, "-r"
	}

	RAW["url"] = []string{url + path}
	RAW["host"] = []string{url}
	RAW["path"] = []string{path}
	RAW["method"] = []string{method}
	RAW["headers"] = lst_raw
	RAW["body"] = []string{body}

	//fmt.Println("::", RAW) //DEBUG
	//os.Exit(0)       //DEBUG - EXIT

	return RAW, ""
}

func Set_LstHeaders(str string) (string, string, bool) {
	/** Split a header string into Header & Value
	 */

	if !strings.Contains(str, ":") {
		return "", "", false
	}

	//Get header & value:
	h := strings.SplitN(str, ":", 2)

	return h[0], h[1], true
}

func Set_ContentType(typ string) (string, string) {
	/** Set content-type
	* Add header value to it's correct post data type
	 */
	var (
		h = "Content-Type"
		v = ""
	)

	switch typ {
	case "post":
		v = "application/x-www-form-urlencoded"
	case "json":
		v = "application/json"
	case "xml":
		v = "application/xml"

	case "none":
		//If no valid type was found it will still proceed because it meight be by purpose by the user:
		//default:
		//fmt.Println("[\033[1;31m FA \033[0m] No valid type was found (post,json,xml) at the end of the post data. This can cause false positives/negatives or crashes.")
	}
	return h, v
}

/** return a random 'User-Agent' from default/selected wordlist in memory*/
func Set_UserAgentRandom() string {
	return G.Lst_RandomAgent[rand.Intn(len(G.Lst_RandomAgent)-1)]
}
