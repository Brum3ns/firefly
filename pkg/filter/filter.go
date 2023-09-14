// Filter is used to filter HTTP responses based on user parse arguments
// In need of the structs : &parse.OptionsConf{}, &storage.Response{}
package filter

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/Brum3ns/firefly/pkg/fail"
	fc "github.com/Brum3ns/firefly/pkg/functions"
	"github.com/Brum3ns/firefly/pkg/option"
	"github.com/Brum3ns/firefly/pkg/request"
)

// FilterData user to filter and match responses
type Filter struct {
	Isset bool
	F     map[string][]string
	M     map[string][]string
	Rules map[string]func(string, float64) bool
	Mode  map[string]string
	Fn    FilterFunctions //Note : Valid = (".", "-",  "++", "--")

	//Local:
	resp request.Response
}

type FilterFunctions struct {
	CompareValue   func(bool, []string, int) bool
	RegexCheck     func(bool, string, string) bool
	HeadersCompare func(bool, *http.Header, string) bool
}

// Filter jobs send to the [chan]nel
type filter_chanJobs struct {
	typ string
	lst []string
}

// Result collectec within the filter process
type filter_chanResult struct {
	status bool
	typ    string
}

// Filtered HTTP responses according to the user configured user options (parse arguments)
func (f Filter) Run(Response request.Response) bool {
	if len(f.F) == 0 && len(f.M) == 0 {
		return false
	}

	//Setup the response that will be analyzed:
	f.resp = Response

	var (
		methodSet      = 2 //(filter|match) = 2 methods max
		mutex          sync.Mutex
		m_methodResult = make(map[string]bool) //Return
	)

	for _, method := range []string{"filter", "match"} {
		var amount, m_methods = f.get_method(method)
		if amount <= 0 {
			methodSet--
			continue
		}
		m_types := make(map[string]bool)

		//setup [chan]nels:
		c_jobs := make(chan filter_chanJobs, amount)
		c_result := make(chan filter_chanResult, amount)
		defer close(c_jobs)
		defer close(c_result)

		//Start 'check' task & Give the jobs: (filter that should be used)
		for t := 0; t < amount; t++ {
			go f.filter(c_jobs, c_result)
		}
		jobs(m_methods, c_jobs)

		//Check result from [chan]nel:
		for i := 0; i < amount; i++ {
			result := <-c_result

			//Calculate amount of "true" hits:
			if result.status {
				mutex.Lock()
				m_types[result.typ] = result.status
				mutex.Unlock()
			}
		}
		//When all checks have been made by the filter/match method. Check if the response did hit the filter in relation to the "mode" in use (and|or):
		m_methodResult[method] = f.modeHit(amount, method, m_types)
	}

	//No filter was set:
	//Note - (Detected in the first *if statement* of the *for loop* above)
	if methodSet == 0 {
		return false
	}

	//Check result and see if the filter for (filter|match) is valid:
	hit := 0
	for method, ok := range m_methodResult {
		//If the method is "match" and its "ok" is false. We need to reverse "ok" to the opposit bool value.
		//Example: "-mc = 200" / "response 200" == "True" and "True" will be threaded as "filtered". Therefore we need to set it to "False" so it won't be filtered (skipped).
		if method == "match" {
			ok = !ok
		}
		//Use double "if" because "hit++" and the last if can be in the same loop:
		if ok {
			hit++
		}
		if hit == len(m_methodResult) {
			return true
		}
	}
	return false
}
func NewFilter(opt *option.Options) *Filter {
	var (
		//Syntax validation for filter rules:
		syntaxValidate = func(ml map[string][]string) bool {
			reSyntax := `^([\d\.]+|(--|\+\+)[\d]+|[\d]+(--|\+\+)|[\d]+-[\d]+)$`
			l_check := []string{"sc", "bs", "lc", "wc", "ha"}
			for _, l := range ml {
				for _, i := range l {

					//Check number based filters:
					if fc.InLst(l_check, i) {

						//Only done in validation process (once) - loop is allowed:
						if ok, _ := regexp.MatchString(reSyntax, i); i != "" && !ok {
							fail.IFFail(1011)
						}
					}
				}
			}
			return true
		}

		isset = func(m map[string][]string) bool {
			for _, lst := range m {
				if len(lst) > 0 && lst[0] != "" {
					return true
				}
			}
			return false
		}

		more = func(InputValue string, compareValue float64) bool {
			v, err := strconv.ParseFloat(InputValue, 64)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return v < compareValue
		}
		less = func(InputValue string, compareValue float64) bool {
			v, err := strconv.ParseFloat(InputValue, 64)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return v > compareValue
		}
		between = func(InputValue string, compareValue float64) bool {
			l := strings.Split(InputValue, "-")
			v1, err1 := strconv.ParseFloat(l[0], 64)
			v2, err2 := strconv.ParseFloat(l[1], 64)
			if err1 != nil || err2 != nil {
				fmt.Println(err1, "AND", err2)
			}
			return v1 < compareValue && v2 > compareValue
		}

		split = func(s string) []string { return strings.Split(s, ",") }
	)

	filter := &Filter{
		Rules: map[string]func(string, float64) bool{
			"++": more,
			"--": less,
			"-":  between,
		},
		Mode: map[string]string{
			"filter": opt.FilterMode,
			"match":  opt.MatchMode,
		},
		F: map[string][]string{
			"sc":  split(opt.FilterCode),
			"bs":  split(opt.FilterSize),
			"lc":  split(opt.FilterLine),
			"wc":  split(opt.FilterWord),
			"ha":  split(opt.FilterHeaderAmount),
			"tm":  split(opt.FilterTime),
			"re":  {opt.FilterBodyRegex},
			"reh": {opt.FilterHeaderRegex},
		},
		M: map[string][]string{
			"sc":  split(opt.MatchCode),
			"bs":  split(opt.MatchSize),
			"lc":  split(opt.MatchLine),
			"wc":  split(opt.MatchWord),
			"ha":  split(opt.MatchHeaderAmount),
			"tm":  split(opt.MatchTime),
			"re":  {opt.MatchBodyRegex},
			"reh": {opt.MatchHeaderRegex},
		},
	}

	//Check if any filter is set and if they have a valid syntax:
	filter.Isset = false
	if isset(filter.F) && syntaxValidate(filter.F) {
		filter.Isset = true
	}
	if isset(filter.M) && syntaxValidate(filter.M) {
		filter.Isset = true
	}

	return filter
}

// Give job to channel runner.MatchData
func jobs(m map[string][]string, c chan<- filter_chanJobs) {
	for k, l := range m {
		c <- filter_chanJobs{
			typ: k,
			lst: l,
		}
	}
}

// Check if the response modeHit the filter in relation to the given user "modeHit"
func (f Filter) modeHit(total int, method string, m map[string]bool) bool {
	if len(m) <= 0 {
		return false
	} else if f.Mode[method] == "or" && len(m) > 0 {
		return true

	} else if (f.Mode[method] == "and") && len(m) == total {
		return true
	}
	return false
}

func (f Filter) get_method(method string) (int, map[string][]string) {
	var ( //Return
		m_methods = make(map[string][]string)
	)
	//Detect what filter method to use:
	var m map[string][]string
	if method == "filter" {
		m = f.F
	} else /*match*/ {
		m = f.M
	}
	//Calculate the amount of filters that are active within the filter method:
	for key, lm := range m {
		if len(lm) > 0 && lm[0] != "" {
			m_methods[key] = lm
		}
	}
	return len(m_methods), m_methods
}

// Validate if the filter/match use a regex (status, size ,word ,line, headerAmount)
func (f Filter) filter(jobs <-chan filter_chanJobs, c chan<- filter_chanResult) {
	var ok bool
	for j := range jobs {
		//[TODO] - (Better code to avoid huge amount of case statements)
		switch j.typ {
		case "sc":
			ok = f.fCompareValue(j.typ, j.lst, float64(f.resp.StatusCode))
		case "bs":
			ok = f.fCompareValue(j.typ, j.lst, float64(f.resp.ContentLength))
		case "wc":
			ok = f.fCompareValue(j.typ, j.lst, float64(f.resp.WordCount))
		case "lc":
			ok = f.fCompareValue(j.typ, j.lst, float64(f.resp.LineCount))
		case "re":
			ok = f.fRegexMatch(j.lst[0], f.resp.Body)
		case "reh":
			ok = f.fRegexMatch(j.lst[0], f.resp.HeaderString)
		case "tm":
			ok = f.fTime(j.typ, j.lst, f.resp.Time)
		case "ha":
			ok = f.fCompareValue(j.typ, j.lst, float64(len(f.resp.Header)))
		}
		c <- filter_chanResult{
			status: ok,
			typ:    j.typ,
		}
	}
}

// Filter response status code AND body size (bs), wordCount (wc), lineCount (lc).
func (f Filter) fTime(typ string, lst []string, respTime float64) bool {
	//Note : (The order of the if condition is important)
	for _, value := range lst {
		valueFloat := getRule(value)

		//Standard
		if v, err := strconv.ParseFloat(value, 64); err == nil && v == respTime {
			return true

			//Advanced (Rules) | Note : (fc.KeepNoneDigits(value) will be equal to a rule: "--", "++", "-")
		} else if fn, ok := f.Rules[valueFloat]; ok && fn(getValue(value), respTime) {
			return true
		}
		continue
	}
	return false
}

// Filter response status code AND body size (bs), wordCount (wc), lineCount (lc).
func (f Filter) fCompareValue(typ string, lst []string, valueResponse float64) bool {
	//Note : (The order of the if condition is important)
	for _, value := range lst {
		//Standard
		if v, err := strconv.ParseFloat(value, 64); err == nil && v == valueResponse {
			return true

			//Advanced (Rules) | Note : (fc.KeepNoneDigits(value) will be equal to a rule: "--", "++", "-")
		} else if fn, ok := f.Rules[getRule(value)]; ok && fn(getValue(value), valueResponse) {
			return true
		}
		continue
	}
	return false
}

// Filter check if given regex within the given string
func (f Filter) fRegexMatch(re string, str string) bool {
	ok, err := regexp.MatchString(re, str)
	fail.IFError("p", err)
	return ok
}

// Remove all none digits (except ".") within a string
func getValue(s string) string {
	var str string
	for _, r := range s {
		if (r >= '0' && r <= '9') || (r == '.' || r == '-') {
			str += string(r)
		}
	}
	str = strings.Replace(str, "--", "", -1)
	return str
}

// Remove all digits within a string
func getRule(s string) string {
	var (
		rule  string
		inRow = 0
	)
	//Note : To avoid looping junk [if/else -> break]
	for _, r := range s {
		if inRow == 2 { //Rule = "++" or "--"
			break
		} else if r == '-' || r == '+' {
			inRow++
			rule += string(r)
		} else if inRow == 1 { //Rule = "-""
			break
		}
	}
	return rule
}
