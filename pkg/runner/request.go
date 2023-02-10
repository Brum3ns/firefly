package runner

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Brum3ns/firefly/pkg/functions"
	fc "github.com/Brum3ns/firefly/pkg/functions"
	"github.com/Brum3ns/firefly/pkg/parse"
	"github.com/Brum3ns/firefly/pkg/storage"
	st "github.com/Brum3ns/firefly/pkg/storage"
)

type client_struct struct {
	client *http.Client
}

//Client configure with custom parse *timeout*:
func Client(opt *parse.Options) *http.Client {
	Timeout := time.Duration(opt.Timeout) * time.Millisecond
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		Timeout:       Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 500,
			MaxConnsPerHost:     500,
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Renegotiation:      tls.RenegotiateOnceAsClient,
			},
		},
	}
	return client
}

//[temp] - all 'functions.InsertData()' will be replaced and improved.
func ReqTemplate(METHOD, URL, PAYLOAD string, opt *parse.Options, verify bool) (*http.Request, string, string, string) {
	/*[Request config]
	*	If the verify is true then a "normal" behavior check is in process, not the fuzz process.
	*	Set payload to a "string" value which can be used as a "int" as well
	*	with 10 as a length since it bypass major most input length. (Password, username, mobile number etc...)
	 */
	var (
		HEADERS = "headers"
		req     *http.Request
		reqErr  error
	)

	//If in verify mode, check if added headers are included to add (in verify request only):
	if verify && opt.VerifyHeader {
		HEADERS = "vheaders"
	}

	//Check Method and if post data should be used
	URL = functions.InsertData(URL, functions.PayloadURLNormalize(PAYLOAD)) //Insert[temp]

	req, reqErr = http.NewRequest(METHOD, URL, nil)

	//Check postData:
	if len(opt.Target["postData"]) > 0 {
		d := functions.InsertData(opt.Target["postData"][0], PAYLOAD) //Insert[temp]

		postData := bytes.NewBuffer([]byte(d))
		req, reqErr = http.NewRequest(METHOD, URL, postData)

		if fc.InLst([]string{"POST", "PUT", "PATCH"}, strings.ToUpper(METHOD)) && len(opt.ReqRaw) <= 0 { //<<<--- START
			h, v := fc.Set_ContentType(strings.ToLower(opt.PostDataType))

			if len(h) > 0 && len(v) > 0 {
				req.Header.Add(h, v)
			}
		}
	}

	//Static headers to be included (static included headers at the bottom of the RAW request):
	if len(opt.Target[HEADERS]) > 0 {
		for _, Header := range opt.Target[HEADERS] {
			Header = functions.InsertData(Header, PAYLOAD) //Insert[temp]

			if h, v, ok := fc.Set_LstHeaders(Header); ok {
				req.Header.Add(h, v)
			}
		}
	}

	//HTTP Raw Headers:
	if len(opt.ReqRaw) > 0 && len(opt.Target["headersRaw"]) > 0 {
		for _, HeaderRaw := range opt.Target["headersRaw"] {
			HeaderRaw = functions.InsertData(HeaderRaw, PAYLOAD) //Insert[temp]

			if h, v, ok := fc.Set_LstHeaders(HeaderRaw); ok {
				req.Header.Add(h, v)
			}
		}
	}

	//Check if the request setup contains errors
	if ok, _ := fc.IFError("", reqErr); ok {
		st.Error++
	}

	if opt.RandomAgent {
		rua := functions.InsertData(fc.Set_UserAgentRandom(), PAYLOAD)
		req.Header.Add("User-Agent", rua)
	} else {
		req.Header.Add("User-Agent", functions.InsertData(opt.UserAgent, PAYLOAD))
	}

	return req, URL, METHOD, PAYLOAD
}

//Reques module that send and add the response data to the "results" channel and use "Response" as struct for dynamic temp variables:
func Request(thread_id int, client *http.Client, opt *parse.Options, jobs <-chan storage.Target, results chan<- storage.Response /*[DELETE?], payloads <-chan storage.Payloads*/) {
	//Extract each urls inside the "job pool":
	for job := range jobs {

		//Setup request template:
		request, url, method, payload := ReqTemplate(job.METHOD, job.URL, job.PAYLOAD, opt, false)

		st.Count++
		respTimer := time.Now()

		//Make request to the target with the given "request template":
		resp, err := client.Do(request)
		if ok, msg := fc.IFError("Check connection, failed to send request.", err); ok {
			results <- storage.Response{
				Error:  true,
				ErrMsg: msg,
			}
			st.Error++
			continue
		}

		//Read the response body to extract data
		if resp.StatusCode > 0 {
			st.Count_valid++

			var bodyBytes, err = ioutil.ReadAll(resp.Body)
			fc.IFError("", err)
			resp.Body.Close() //FIX

			contentLength := int64(len(string(bodyBytes[:])))
			URLNoPayload := strings.ReplaceAll(job.URL, opt.Insert, "")

			//Store the response data into the "results" chan
			results <- storage.Response{
				ThreadID: thread_id,
				Id:       job.ID,

				Tag:           job.TAG,
				Payload:       payload,
				Valid:         st.Count_valid,
				Url:           url,
				UrlNoPayload:  URLNoPayload,
				Method:        method,
				Status:        resp.StatusCode,
				Body:          bodyBytes,
				Headers:       resp.Header,
				HeadersString: fc.HeadersToStr(resp.Header),
				ContentLength: contentLength,
				ContentType:   strings.Replace(resp.Header.Get("Content-Type"), "; charset=UTF-8", "", 1),
				RespTime:      float64(time.Since(respTimer).Seconds()),
			}

		}

		//Time delay:
		time.Sleep(time.Duration(opt.Delay) * time.Duration(time.Millisecond))
	}
}
