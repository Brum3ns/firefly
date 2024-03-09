package output

import (
	"net/http"

	"github.com/Brum3ns/firefly/pkg/difference"
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/transformation"
)

// The output result is the final result that is generated that stores all the result from all processes done by the runner
// Note : (The final result only stores the result details that are of int `json:""`erest to the user and not the properties that were used during the runner process. The variable may be reformulated for better readability)
type ResultFinal struct {
	RequestId      int      `json:"RequestId"`
	TargetHashId   string   `json:"TargetId"`
	Tag            string   `json:"Tag"`
	Date           string   `json:"Date"`
	Payload        string   `json:"Payload"`
	Request        Request  `json:"Request"`
	Response       Response `json:"Response"`
	Scanner        Scanner  `json:"Scanner"`
	Error          error    `json:"Error"`
	OK             bool     `json:"-"`
	UnkownBehavior bool
	//Origin       string   `json:"Origin"`
	//Behavior     Behavior `json:"Behavior"`
}

//TODO
// Contains the collected Behavior from the target
/* type Behavior struct {
	CWE         string `json:"CWE"`
	Component   string `json:"Component"`
	Confidence  string `json:"Confidence"`
	Description string `json:"Desc"`
} */

// Refer to the results of the Request/response process
type Request struct {
	URL         string      `json:"URL"`
	URLOriginal string      `json:"URL-Original"`
	Host        string      `json:"Host"`
	Scheme      string      `json:"Scheme"`
	Method      string      `json:"Method"`
	PostBody    string      `json:"PostBody"`
	Proto       string      `json:"HTTP"`
	Headers     [][2]string `json:"Headers"`
}

// Refer to the results of the request/Response process
type Response struct {
	StatusCode    int         `json:"Status-Code"`
	WordCount     int         `json:"WordCount"`
	LineCount     int         `json:"LineCount"`
	HeaderAmount  int         `json:"Header-Amount"`
	ContentLength int         `json:"Content-Length"`
	ContentType   string      `json:"Contnet-Type"`
	Host          string      `json:"Host"`
	Body          string      `json:"-"`
	Title         string      `json:"Title"`
	Proto         string      `json:"HTTP"`
	IPAddress     []string    `json:"IPAddress"`
	Time          float64     `json:"Response-Time"`
	Headers       http.Header `json:"Headers"`
}

// Refer to the results of the scanning process
type Scanner struct {
	Extract        extract.Result        `json:"Extract"`
	Diff           difference.Result     `json:"Diff"`
	Transformation transformation.Result `json:"Transformation"`
}
