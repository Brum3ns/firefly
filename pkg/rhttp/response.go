package rhttp

import (
	"io"
	"math"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type Response struct {
	Status        string      `json:"status"`
	StatusCode    int         `json:"statuscode"`
	WordCount     int         `json:"wordcount"`
	LineCount     int         `json:"linecount"`
	HeaderAmount  int         `json:"headeramount"`
	BodySize      int         `json:"bodysize"`
	ContentLength int64       `json:"contentlength"`
	ContentType   string      `json:"contenttype"`
	Host          string      `json:"host"`
	Body          string      `json:"body"`
	Proto         string      `json:"proto"`
	Time          float64     `json:"time"`
	Headers       http.Header `json:"headers"`
	//*http.Response
}

func NewResponse(resp *http.Response, respTime time.Duration) (Response, error) {
	bodyBytes, err := httputil.DumpResponse(resp, true)

	//bodyBytes, err := MakeResponseBody(resp.Body)
	if err != nil {
		return Response{}, err
	}
	bodyStr := string(bodyBytes)

	//ipAddress, _ := net.LookupIP(hostname)

	return Response{
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		WordCount:     GetBodyWordCount(bodyStr),
		LineCount:     GetBodyLineCount(bodyStr),
		HeaderAmount:  GetHeaderAmount(resp.Header),
		BodySize:      GetBodySize(bodyStr),
		ContentLength: resp.ContentLength,
		ContentType:   GetContentType(resp.Header),
		//Host:          hostname,
		Body:    bodyStr,
		Proto:   resp.Proto,
		Time:    MakeResponseTime(respTime),
		Headers: resp.Header,
	}, nil
}

func MakeResponseBody(body io.Reader) ([]byte, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return []byte{}, err
	}
	return bodyBytes, nil
}

func MakeResponseTime(respTime time.Duration) float64 {
	return math.Round(respTime.Seconds()*10000) / 10000
}

func GetContentType(headers http.Header) string {
	if ct, ok := headers["content-type"]; ok && len(ct) > 0 {
		return ct[0]
	}
	return ""
}

func GetBodySize(body string) int {
	return len(body)
}

func GetBodyLineCount(body string) int {
	return len(strings.Split(body, "\n"))
}

func GetBodyWordCount(body string) int {
	return len(strings.Fields(body))
}

func GetHeaderAmount(headers http.Header) int {
	count := 0
	for _, values := range headers {
		count += len(values)
	}
	return count
}
