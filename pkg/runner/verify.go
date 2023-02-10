/**NOTE: This process should not run in a gorutine.*/
package runner

import (
	"regexp"
	"strings"

	"github.com/Brum3ns/FireFly/pkg/firefly/prepare"
	"github.com/Brum3ns/FireFly/pkg/functions"
	fc "github.com/Brum3ns/FireFly/pkg/functions"
	G "github.com/Brum3ns/FireFly/pkg/functions/globalVariables"
	"github.com/Brum3ns/FireFly/pkg/storage"
)

type verifyResult struct {
	s_typ    string
	s_regx   string
	l_bodies []string

	s_payload        string
	s_payloadMark    string
	s_payloadPattern string
	i_payloadReflect int
	m_HTMLAttributes map[string][]int

	m_regexDynamic    map[string][]string
	m_patternsDynamic map[string][]string
	m_linesDynamic    map[string][]string

	//Temporary (Used to calculate unique lines/words)
	tmpM_lineCount map[string]int
	tmpM_wordCount map[string]int

	//Return
}

func (r *Runner) Verify() *storage.VerifyData {
	r.verifyData.Target = r.verifier.Data

	for _, vd := range r.verifyData.Target {
		if vd.Tag == "verifyPayload" {
			data := prepare.Body(vd.Body, vd.Payload)
			r.verifyData.TargetInfo[vd.URL] = append(r.verifyData.TargetInfo[vd.URL], data)
		}
	}

	r.verifyData.DyncRegex = VerifyDynamic(*r)
	r.verifyData.CharEncode = VerifyChars(*r)

	return r.verifyData
}

/* func VerifyData(r Runner) *verifyTargetData {

} */

//[TODO] - Move to functions (Improve)
func VerifyDynamic(r Runner) map[string][]string {
	/**Extract unique dynamic gadget from verified responses to be used for diff analyze in future processes,
	* Return two maps containing the VResp dynamic data for -> lines & words.
	 */

	data := &verifyResult{
		s_typ:            "body", //[TODO] Headers as well
		s_payloadMark:    "{PAYLOAD}",
		s_payloadPattern: G.PayloadPattern,
		i_payloadReflect: 0,

		s_regx: functions.PayloadRegexMark(1),

		m_HTMLAttributes:  make(map[string][]int),
		m_regexDynamic:    make(map[string][]string),
		m_patternsDynamic: make(map[string][]string),
		m_linesDynamic:    make(map[string][]string),

		tmpM_lineCount: make(map[string]int),
		tmpM_wordCount: make(map[string]int),
	}

	for _, body1 := range data.l_bodies {
		l_body1 := strings.Split(body1, "\n")
		for _, body2 := range data.l_bodies {
			//l_body2 := strings.Split(body2, "\n")

			if body1 != body2 {
				//Extract each line and add non duplicates related to 'x' body list of lines:
				for _, line := range l_body1 {
					//if !fc.InLst(l_body2, line) {
					data.tmpM_lineCount[line] += 1

					//Extract each word and add non duplicates related to 'x' body list of lines:
					for _, word := range fc.ToLstSplitRe(line, "[^a-zA-Z0-9-_.:]") {
						data.tmpM_wordCount[word] += 1
					}
					//}
				}
			}
		}
	}

	for uniqueLine, nr := range data.tmpM_lineCount {
		// The amount of request verification will always repeat same word: itself -1 (if atleast 2>). Therefore we push it down to 1)
		if G.Verify > 1 && nr == G.Verify-1 {
			//Dynamic line finished. Add it to map:
			data.m_linesDynamic[data.s_typ] = append(data.m_linesDynamic[data.s_typ], uniqueLine)
		}
	}
	for uniqueWord, nr := range data.tmpM_wordCount {
		if G.Verify > 1 && nr == G.Verify-1 {
			data.m_patternsDynamic[data.s_typ] = append(data.m_patternsDynamic[data.s_typ], uniqueWord)
		}
	}

	data.m_regexDynamic[data.s_typ] = r.DyncGetRegexPtn(data.s_typ, data.m_linesDynamic, data.m_patternsDynamic)

	return data.m_regexDynamic
}

func (r *Runner) DyncGetRegexPtn(t string, mL, mP map[string][]string) []string {
	var (
		l = []string{}
		//vPayload = (r.options.PayloadPattern + (r.options.VerifyPayload) + r.options.PayloadPattern)
		tmpMark = "__FIREFLY_DYNC__"
	)

	for _, line := range mL[t] {
		for _, ptn := range mP[t] {
			if strings.Contains(line, ptn) {

				line = strings.ReplaceAll(line, ptn, tmpMark)
				//line = strings.ReplaceAll(line, G.PayloadMark, tmpMark)

				//[Temp][TODO] - Fix more fancy & faster
				for strings.Contains(line, (tmpMark + tmpMark)) {
					line = strings.ReplaceAll(line, tmpMark, ptn)
				}

			}
		}
		//Escape to a valid regex string & add regex pattern to map:
		//line = ReEscp.Replace(line)
		line = functions.RuneToASCII(line)
		line = regexp.QuoteMeta(line)

		rePtn := (strings.ReplaceAll(line, tmpMark, `*(.*?)*`))
		//rePtn := ("^" + strings.ReplaceAll(line, tmpMark, `*(.*?)*`) + "$")

		if !fc.InLst(l, rePtn) {
			l = append(l, rePtn)
		}
	}

	return l
}

func VerifyChars(r Runner) map[string]map[rune][]string {
	var (
		//vR = r.vresp.Data //vRes.Data
		//re = regexp.MustCompile(`[^0-9]+`)
		t = "body"
	)
	m := make(map[string]map[rune][]string)
	m[t] = make(map[rune][]string)

	for _, data := range r.verifyData.Target {
		body := string(data.Body[:])

		for rn, i := range r.wordlist.VerifyPayloadChars {
			re := strings.Replace(i, (string(rn)), `(.*?)`, 1)

			if l := fc.RegexBetweenLst(body, re); len(l) > 0 && i != r.options.VerifyPayload {
				m[t][rn] = l
			}
		}
	}

	return m
}
