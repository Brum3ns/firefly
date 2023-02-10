/**This types should only be used as templates*/
package types

//[TODO]
/* type VerifyTarget struct {
	Tag           string
	URL           string
	Payload       string
	Body          string
	HeadersString string
	Headers       http.Header
} */

/**pointer to struct 'prepare.BodyData'*/
type VerifyInfoData struct {
	BodyPrepared string

	PayloadRegex        string
	PayloadMark         string
	PayloadReflectCount int

	LineCount          int
	WordCount          int
	HTMLTagCount       int
	HTMLAttrCount      int
	HTMLAttrValueCount int

	ItemLst []string

	Lines         map[string]int
	Words         map[string]int
	HTMLTags      map[string]int
	HTMLAttr      map[string]int
	HTMLAttrValue map[string]int
	HTML          map[string][]string
}

func Setup_VerifyInfoData() *VerifyInfoData {
	data := &VerifyInfoData{
		ItemLst:       []string{"word", "line", "htmltag", "htmlattr", "htmlattrvalue", "payloadreflect"},
		Lines:         make(map[string]int),
		Words:         make(map[string]int),
		HTMLTags:      make(map[string]int),
		HTMLAttr:      make(map[string]int),
		HTMLAttrValue: make(map[string]int),
		HTML:          make(map[string][]string),
	}

	return data
}
