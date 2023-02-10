package storage

import (
	"net/http"

	"github.com/Brum3ns/firefly/pkg/firefly/types"
)

//[TODO] Move to struct in parse.Options + yml support
var (
	StaticFilePath = "db/resources/detection/"
)

type Target struct {
	ID int

	TAG     string
	URL     string
	METHOD  string
	PAYLOAD string

	RAW []string
}

//Primary data storage:
type Collection struct {
	ThreadID     int
	ID           int
	Tag          string
	Payload      string
	PayloadClear string

	Url          string
	UrlNoPayload string
	Method       string
	StatusCode   int

	ContentLength int64
	ContentType   string
	//HeadersString string
	Headers http.Header

	RespTime float64

	Body         []byte
	BodyByteSize int

	LineCount int
	WordCount int

	//[Diff]erence detected:
	AvgAmountDiff map[string]int
	//RespDiff      map[string]map[string]string //Target -> Payload, Amount, Diff(s)

	//Trace Error detected:
	RespErr       map[string]map[int][]string //Type, Amount, Error(s)
	RespErrAmount map[string]int

	//Transformations:
	Transformation map[string]string //Payload, TransformedTo
	Tfmt_display   string
	Tfmt_ok        bool

	//Filter
	Filter map[string]map[string]bool //Type(F/M), FilterType, status (T/F)

	//Technologies
	//BakendFunc   map[int][]string

	Valid        int
	HeadersMatch []string
	RegexMatch   bool

	//TaskStatus map[int][]string //RespID, taskNames (success/fail) - (Should always be 4 in the list)

	//Verbose:
	Icon   string
	HasErr string

	//Error debug:
	ErrMsg        string
	Errors        bool
	Skip          bool
	TaskStatus    map[string]map[int]bool
	VerifyProcess bool
	Status        bool
}

type Response struct {
	ThreadID int
	Tag      string
	Id       int

	Payload string

	Url           string
	UrlNoPayload  string
	Method        string
	ContentLength int64
	ContentType   string

	HeadersString string
	Headers       http.Header

	Body []byte

	RespTime     float64
	BodyByteSize int
	Status       int
	Valid        int

	Skip bool

	//Error debug:
	Error  bool
	ErrMsg string
}

type Payloads struct {
	Payload    string
	PayloadYML string
}

type VerifyData struct {
	Target map[int]struct { //[TODO] v1.1 - Setup a 'types.go' file to have all the struct templates
		Tag           string
		URL           string
		Payload       string
		Body          string
		HeadersString string
		Headers       http.Header
	}

	TargetInfo map[string][]types.VerifyInfoData
	/*struct { //[TODO] v1.1 Setup a 'types.go' file to have all the struct templates
		 bodyPrepared string

		regx_payload        string
		payloadMark         string
		payloadReflectCount int

		lineCount     int
		wordCount     int
		lines         map[string]int
		words         map[string]int
		HTMLTags      map[string]int
		HTMLAttr      map[string]int
		HTMLAttrValue map[string]int
		HTML          map[string][]string
	} */

	//type (body, header etc..), values:
	DyncWord  map[string][]string
	DyncLine  map[string][]string
	DyncRegex map[string][]string

	//Payload behavior
	CharEncode    map[string]map[rune][]string // type, char, char-form
	PayloadModify map[string][]string          //Old-Payload, New-Payload

	VR_AvgRespTime      []float64
	FalsePositiveErrors []string
	LstG_Tag            map[string]int
}

type Wordlists struct {
	UseTechniques    []string
	Valid_Techniques []string

	//Auto detected param[s] from user input (GET/POST)
	Lst_paramsURL  []string
	Lst_paramsData []string

	//all types of patters, errors, keywords to look for:
	MG_patterns map[string][]string //type, lstOfPatterns - total lists [16]

	//Lists payloads [8]
	VerifyPayloadChars     map[rune]string
	TransformationCompare  map[string]string
	TransformationPayloads []string
	Lst_default            []string
	Fuzz                   []string
	Reflect                []string
	CachePoisoning         []string
	Headers                []string
	Directories            []string
	Timebased              []string
}

func ConfWordlist() *Wordlists {
	wl := &Wordlists{}

	wl.MG_patterns = make(map[string][]string)
	wl.Valid_Techniques = []string{"fuzz", "transformation", "reflect", "cp", "headers", "dir", "time"}
	return wl
}

/**Config - (Response struct setup)*/
func ConfResp() *Response {
	return &Response{}
}

func ConfVerifyData() *VerifyData {
	vresp := &VerifyData{}
	vresp.TargetInfo = make(map[string][]types.VerifyInfoData)
	vresp.Target = make(map[int]struct {
		Tag           string
		URL           string
		Payload       string
		Body          string
		HeadersString string
		Headers       http.Header
	})
	vresp.LstG_Tag = map[string]int{ //[TODO]
		"cache":           0,
		"cms":             0,
		"dbms":            0,
		"errors":          0,
		"errorMariaDB":    0,
		"errorPHP":        0,
		"errorPostgreSQL": 0,
		"errorPython":     0,
		"extension":       0,
		"filename":        0,
		"infoleak":        0,
		"pattern":         0,
		"patternJs":       0,
		"technology":      0,
		"webservice":      0,
		"other":           0,
	}
	return vresp
}
