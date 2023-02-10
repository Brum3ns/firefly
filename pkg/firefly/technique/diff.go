package technique

import (
	"fmt"

	"github.com/Brum3ns/FireFly/pkg/firefly/prepare"
	"github.com/Brum3ns/FireFly/pkg/firefly/types"
	fc "github.com/Brum3ns/FireFly/pkg/functions"
	"github.com/Brum3ns/FireFly/pkg/storage"
)

type diffData struct {
	URL     string
	payload string

	info   types.VerifyInfoData
	resp   storage.Response
	verify *storage.VerifyData

	AvgWordAmount          []int
	AvgLineAmount          []int
	AvgHTMLAmount          []int
	AvgHTMLAttrAmount      []int
	AvgHTMLAttrValueAmount []int

	AvgSumDiff map[string]int

	//Return
	resultBanner string
	result       map[string]string
}

//"body", resp, r.verifyData
func Diff(typ string, response storage.Response, verifyData *storage.VerifyData) (map[string]string, map[string]int) {
	data := &diffData{
		URL:     response.UrlNoPayload,
		payload: response.Payload,

		resp:   response,
		verify: verifyData,
		info:   *types.Setup_VerifyInfoData(),

		AvgSumDiff: make(map[string]int),
		result:     make(map[string]string),
	}
	//Setup fuzzed response data information to be compared with stored verifiyed data information:
	data.info = prepare.Body(string(data.resp.Body[:]), data.payload)

	//If response data do not exist:
	if len(data.verify.Target[data.resp.Id].URL) <= 0 {
		return nil, nil
	}

	//Extract known information about the fuzzed target:
	m_AvgAmount := make(map[string][]int)
	l_keyItem := data.info.ItemLst
	for _, vfInfo := range data.verify.TargetInfo[data.URL] {

		/*Compare if the body if the fuzzed response is equal to the stored verified body.
		Major of case it's a diff since web application almost always contains some sort of dynamic content*/
		if data.info.BodyPrepared != vfInfo.BodyPrepared {

			/*Be aware the order do matter
			l_keyItem = ("word", "line", "htmltag", "htmlattr", "htmlattrvalue", "payloadreflect")*/

			//[Words]
			if v := data.DiffAmount(data.info.WordCount, vfInfo.WordCount); true {
				id := 0
				m_AvgAmount[l_keyItem[id]] = append(m_AvgAmount[l_keyItem[id]], v)
			}
			//[Lines]
			if v := data.DiffAmount(data.info.LineCount, vfInfo.LineCount); true {
				id := 1
				m_AvgAmount[l_keyItem[id]] = append(m_AvgAmount[l_keyItem[id]], v)
			}
			//[HTMLTag]
			if v := data.DiffAmount(data.info.HTMLTagCount, vfInfo.HTMLTagCount); true {
				id := 2
				m_AvgAmount[l_keyItem[id]] = append(m_AvgAmount[l_keyItem[id]], v)
			}
			//[HTMLAttr]
			if v := data.DiffAmount(data.info.HTMLAttrCount, vfInfo.HTMLAttrCount); true {
				id := 3
				m_AvgAmount[l_keyItem[id]] = append(m_AvgAmount[l_keyItem[id]], v)
			}
			//[HTMLAttrValue]
			if v := data.DiffAmount(data.info.HTMLAttrValueCount, vfInfo.HTMLAttrValueCount); true {
				id := 4
				m_AvgAmount[l_keyItem[id]] = append(m_AvgAmount[l_keyItem[id]], v)
			}
			//[Payloadreflect]
			if v := data.DiffAmount(data.info.PayloadReflectCount, vfInfo.PayloadReflectCount); true {
				id := 5
				m_AvgAmount[l_keyItem[id]] = append(m_AvgAmount[l_keyItem[id]], v)
			}

			//Analyze difference in deep:
			//data.AvgWordAmount = append(data.AvgWordAmount, data.DiffAmount(data.info.WordCount, vfInfo.WordCount))
			//data.AvgLineAmount = append(data.AvgLineAmount, data.DiffAmount(data.info.LineCount, vfInfo.LineCount))
		}

	}

	//[TODO] - Improve code below structure:
	for _, keyItem := range l_keyItem {
		data.DiffAvgSum(keyItem, m_AvgAmount[keyItem])
	}

	return data.result, data.AvgSumDiff
}

/**Recive data from fuzzed response and the stored response to analyze their data differences*/
func (d diffData) DiffAnalyze() {
	for _, vfData := range d.verify.TargetInfo[d.URL] {

		fmt.Println(vfData.LineCount)
	}
}

func (d diffData) DiffAmount(count, vfCount int) int {
	return fc.LengthDiff(count, vfCount)

}

/* func (d diffData) DiffAvgAmount(lst []int) int {
	return fc.LstSum(lst) / len(lst)
}
*/
func (d diffData) DiffAvgSum(item string, lst []int) {
	/*Setup average sum value that difference from 'info.ItemLst' which equals:
	=> ("word", "line", "htmltag", "htmlattr", "htmlattrvalue", "payloadreflect")*/

	if sum := fc.LstSum(lst); sum > 0 && len(lst) > 0 {
		d.AvgSumDiff[item] = (fc.LstSum(lst) / len(lst))
	} else {
		d.AvgSumDiff[item] = 0
	}
}
