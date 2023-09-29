package request

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/firefly/global"
	"github.com/Brum3ns/firefly/pkg/functions"
	fc "github.com/Brum3ns/firefly/pkg/functions"
	"github.com/Brum3ns/firefly/pkg/random"
)

var regex_HTMLTitle = regexp.MustCompile(`<title>(.*?)<\/title>`)

type Result struct {
	RequestId    int
	TargetHashId string
	Tag          string
	Date         string
	Payload      string
	Request      Request
	Response     Response
	Skip         bool
	Error        error
}

type Response struct {
	Time            float64
	WordCount       int
	LineCount       int
	HeaderAmount    int
	ContentLength   int
	ContentType     string
	Title           string
	Body            string
	HeaderString    string
	IPAddress       []string
	HeadersOriginal [][2]string
	http.Response
}

// Request configuration (alias of the "http.Request" struct but with some extra variables added)
type Request struct {
	Body            string
	URLOriginal     string
	HeadersOriginal [][2]string
	http.Request
}

type RequestProperties struct {
	RequestId       int
	RandomUserAgent bool
	TargetHashId    string
	Tag             string
	URL             string
	URLOriginal     string
	Payload         string
	PostBody        string
	Method          string
	HeadersOriginal [][2]string
	Headers         http.Header
}

type ClientProperties struct {
	Timeout             int
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	HTTP2               bool
	Proxy               string
}

var (
	regexScheme      = regexp.MustCompile(`^*(.*?)*://`)
	randomUserAgents = fc.FileToList(global.FILE_RANDOMAGENT)
	ruaLength        = len(randomUserAgents)
)

// Reques module that send and add the response data to the "results" channel and use "Response" as struct for dynamic temp variables:
func (w worker) request(requestProperties RequestProperties) Result {

	httpRequest, err := http.NewRequest(requestProperties.Method, requestProperties.URL, SetPostbody(requestProperties.PostBody))
	if err != nil {
		return Result{Error: err}
	}

	//Add headers:
	httpRequest.Header = requestProperties.Headers

	if ruaLength > 0 && requestProperties.RandomUserAgent {
		httpRequest.Header.Add("User-Agent", randomUserAgents[random.Rand.Intn(ruaLength)])
	}

	var responseTime float64
	Timer := time.Now()
	response, err := w.client.Do(httpRequest)
	if err != nil {
		return Result{Error: err}
	}
	//The response was successful. Get the response time:
	if len(response.Status) > 0 {
		responseTime = float64(time.Since(Timer).Seconds())
	}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(httpRequest.Body)

	//Read the response body content:
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(design.STATUS.ERROR, "Could not read the response body:", err)
		return Result{Error: err}
	}

	bodyString := string(bodyBytes[:])
	response.Body.Close()

	//In case any normalization happens within the request post body it will be spotted for the userlater:
	return Result{
		TargetHashId: requestProperties.TargetHashId,
		RequestId:    requestProperties.RequestId,
		Tag:          requestProperties.Tag,
		Payload:      requestProperties.Payload,
		Date:         time.Now().Format(time.UnixDate),
		Error:        nil,
		Request: Request{
			URLOriginal:     requestProperties.URLOriginal,
			HeadersOriginal: requestProperties.HeadersOriginal,
			Request:         *httpRequest,
		},
		Response: Response{
			IPAddress:     GetIPAddresses(response.Request.URL.Hostname()),
			HeaderString:  headersToStr(response.Header),
			Title:         GetHTMLTitle(bodyString),
			ContentType:   response.Header.Get("content-type"),
			ContentLength: len(bodyString),
			HeaderAmount:  len(response.Header),
			Time:          responseTime,
			LineCount:     functions.LineCount(bodyString),
			WordCount:     functions.WordCount(bodyString),
			Body:          bodyString,
			Response:      *response,
		},
	}
}

func URLNormalize(s string) string {
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

// Client configure with custom parse *timeout*:
func setClient(p *ClientProperties) *http.Client {
	var (
		proxy   = http.ProxyFromEnvironment
		timeout = time.Duration(time.Duration(p.Timeout) * time.Second)
	)
	if len(p.Proxy) > 0 {
		if p, err := url.Parse(p.Proxy); err == nil {
			proxy = http.ProxyURL(p)
		}
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		Timeout:       timeout,
		Transport: &http.Transport{
			ForceAttemptHTTP2:   p.HTTP2,
			Proxy:               proxy,
			MaxIdleConns:        p.MaxIdleConns,
			MaxIdleConnsPerHost: p.MaxIdleConnsPerHost,
			MaxConnsPerHost:     p.MaxConnsPerHost,
			DialContext: (&net.Dialer{
				Timeout: timeout,
			}).DialContext,
			TLSHandshakeTimeout: timeout,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS10,
				Renegotiation:      tls.RenegotiateOnceAsClient,
			},
		},
	}
	return client
}

// Get a list of IP addresses that the hostname resolves to
func GetIPAddresses(hostname string) []string {
	var lst []string
	ips, _ := net.LookupIP(hostname)
	for _, i := range ips {
		lst = append(lst, i.String())
	}
	return lst
}

// Setup the post data used within the request (if any)
func SetPostbody(body string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(body))
}

// Check if a URL contains a scheme (any type)
// Return the scheme or an empty string if no scheme was presented in the given URL.
func ContainScheme(s string) string {
	if lst_scheme := regexScheme.FindStringSubmatch(s); len(lst_scheme) > 0 {
		return lst_scheme[1]
	}
	return ""
}

// Validate if the scheme is either http or https
func ValidScheme(s string) bool {
	return s == "http" || s == "https"
}

func ValidURLOrIP(s string) bool {
	_, err := url.Parse(s)
	return err == nil || net.ParseIP(s) != nil
}

// return a random 'User-Agent' from default/selected wordlist in memory
func randomUserAgent() string {
	return randomUserAgents[rand.Intn(len(global.RANDOM_AGENTS)-1)]
}

// Convert the *http.Header* to a string (type: "map[string][]string").
// The converted string version is sorted which makes it easier to compare with others.
func headersToStr(headers http.Header) string {
	if headers == nil {
		return ""
	}

	arr := make([]string, 0, len(headers))
	for k, v := range headers {
		arr = append(arr, fmt.Sprintf("%s: %s\n", k, strings.Join(v, " ")))
	}
	sort.Strings(arr)
	return strings.Join(arr, "")
}

// Use regexp to extract the HTML title and return it as a string
func GetHTMLTitle(s string) string {
	var title string
	if ti := regex_HTMLTitle.FindString(s); ti != "" {
		title = ti[7 : len(ti)-8] //(Known size from re_title)
	}
	return title
}
