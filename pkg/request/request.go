package request

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
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

	"github.com/Brum3ns/firefly/pkg/parameter"
)

var regex_HTMLTitle = regexp.MustCompile(`<title>(.*?)<\/title>`)

type Result struct {
	RequestId    int
	TargetHashId string
	Tag          string
	Date         string
	Payload      string
	Request      HttpRequest
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

// HttpRequest configuration (alias of the "http.HttpRequest" struct but with some extra variables added)
type HttpRequest struct {
	Body            string
	URLOriginal     string
	HeadersOriginal [][2]string
	http.Request
}

// Request settings for each individuallyrequest
type RequestSettings struct {
	RequestId    int
	TargetHashId string
	Tag          string
	URL          string
	URLOriginal  string
	Payload      string
	Method       string
	UserAgents   []string
	Parameter    parameter.Parameter
	RequestBase
}

// Stores the base (static) HTTP data that will be used within all requests
type RequestBase struct {
	PostBody             string
	InsertPoint          string
	RandomUserAgent      bool
	HeadersOriginalArray [][2]string
	Headers              http.Header
}

type ClientSettings struct {
	Timeout             int
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	HTTP2               bool
	Proxy               string
}

type Host struct {
	URL    string
	Scheme string
	Method string
}

var (
	regexScheme = regexp.MustCompile(`^(.*?)://`)
)

// Reques module that send and add the response data to the "results" channel and use "Response" as struct for dynamic temp variables:
func Request(client *http.Client, requestSettings RequestSettings) Result {
	httpRequest, err := http.NewRequest(requestSettings.Method, requestSettings.URL, SetPostbody(requestSettings.PostBody))
	if err != nil {
		return Result{Error: err}
	}

	//Add headers:
	httpRequest.Header = requestSettings.Headers

	//Add random headers (if set):
	ruaLength := len(requestSettings.UserAgents)
	if ruaLength > 0 && requestSettings.RandomUserAgent {
		httpRequest.Header.Add("User-Agent", requestSettings.UserAgents[rand.Intn(ruaLength-1)])
	}

	Timer := time.Now()
	response, err := client.Do(httpRequest)
	if err != nil {
		return Result{Error: err}
	}
	//The response was successful. Get the response time:
	var responseTime float64
	if len(response.Status) > 0 {
		responseTime = float64(time.Since(Timer).Seconds())
	}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(httpRequest.Body)

	//Read the response body content:
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Could not read the response body:", err)
		return Result{Error: err}
	}

	bodyString := string(bodyBytes[:])
	response.Body.Close()

	//In case any normalization happens within the request post body it will be spotted for the userlater:
	return Result{
		TargetHashId: requestSettings.TargetHashId,
		RequestId:    requestSettings.RequestId,
		Tag:          requestSettings.Tag,
		Payload:      requestSettings.Payload,
		Date:         time.Now().Format(time.UnixDate),
		Error:        nil,
		Request: HttpRequest{
			URLOriginal:     requestSettings.URLOriginal,
			HeadersOriginal: requestSettings.HeadersOriginalArray,
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
			LineCount:     len(strings.Split(bodyString, "\n")),
			WordCount:     len(strings.Fields(bodyString)),
			Body:          bodyString,
			Response:      *response,
		},
	}
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

// Client configure with custom parse *timeout*:
func NewClient(p ClientSettings) *http.Client {
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

// Setup the post data used within the request (if any)
func SetPostbody(body string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(body))
}

// Check if a URL contains a scheme (any type)
// Return the scheme or an empty string if no scheme was presented in the given URL.
func ContainScheme(s string) string {
	fmt.Println(s)
	if lst_scheme := regexScheme.FindStringSubmatch(s); len(lst_scheme) > 0 {
		fmt.Println(lst_scheme)
		return lst_scheme[1]
	}
	return ""
}

// Validate if the scheme is either http or https
func ValidHTTPScheme(s string) bool {
	return s == "http" || s == "https"
}

func ValidURLOrIP(s string) bool {
	_, err := url.Parse(s)
	return err == nil || net.ParseIP(s) != nil
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

// Set a new value for an existing header
func SetNewHeaderValue(arr [][2]string, header string, value string) [][2]string {
	header = strings.ToLower(header)
	for idx, h := range arr {
		if strings.ToLower(h[0]) == header {
			arr[idx] = [2]string{h[0], value}
			break
		}
	}
	return arr
}

// Get a header and it's value from a header array list ([][2]string)
// Note : (Requested header name  is in-case sensitive)
func GetHeader(arr [][2]string, header string) (string, string) {
	header = strings.ToLower(header)
	for _, h := range arr {
		if strings.ToLower(h[0]) == header {
			return h[0], h[1] //header, value
		}
	}
	return "", ""
}

// Use regexp to extract the HTML title and return it as a string
func GetHTMLTitle(s string) string {
	var title string
	if ti := regex_HTMLTitle.FindString(s); ti != "" {
		title = ti[7 : len(ti)-8] //(Known size from re_title)
	}
	return title
}

// Make a unique md5 hash from the url and method:
func MakeHash(Url, method string) string {
	hash := md5.Sum([]byte(method + Url))
	return hex.EncodeToString(hash[:])
}

// Take a full raw URL and return the raw query
func GetRawQuery(Url string) (string, error) {
	u, err := url.Parse(Url)
	return u.RawQuery, err
}
