package technique

import (
	"fmt"
	"strings"

	"github.com/Brum3ns/FireFly/pkg/design"
	fc "github.com/Brum3ns/FireFly/pkg/functions"
	"github.com/Brum3ns/FireFly/pkg/storage"
)

type Edata struct {
	lErr  []string
	tErr  string
	count int
}

/**[TODO] - Will be imporved a lot*/
func ErrorDetect(body string, wl *storage.Wordlists, verifyData *storage.VerifyData) (map[string]map[int][]string, map[string]int, string, bool) {
	/** Detect errors/keywords in the a HTTP response body
	* Return the amount(count) of errors(lerr) and it's type(errType)
	 */
	var (
		loops     = len(wl.MG_patterns)
		ErrBanner string
		ErrCount  int
		OK        bool
	)
	m_errAmount := make(map[string]int)

	m := make(map[string]map[int][]string)
	c := make(chan Edata, loops)
	defer close(c)

	//Start background process (mini engine)
	go Extract(body, wl, verifyData.FalsePositiveErrors, c)

	//Listener:
	for i := 0; i < loops; i++ {
		results := <-c
		ErrCount += results.count
		m[results.tErr] = make(map[int][]string)
		m[results.tErr][results.count] = results.lErr

		for _, item := range results.lErr {
			m_errAmount[item] += 1
		}
	}

	//Banner & Icon design for verbose (Detected atleast one error/pattern in response):
	if ErrCount > 0 {
		ErrBanner = (design.RespError + fmt.Sprintf("(\033[31m%v\033[0m)", ErrCount))
		OK = true
	} else {
		ErrBanner = design.Null
	}

	return m, m_errAmount, ErrBanner, OK
}

func Extract(body string, wl *storage.Wordlists, vres_lFN []string, result chan<- Edata) {
	/** Child of function 'ErrorDetect()'
	 *	Return the result of amount, errorType, errors
	 */

	//Extract all the wordlist and types from wl.MGpatterns and use them to collect patterns:
	for tkey, l := range wl.MG_patterns {
		count, lst := 0, []string{}

		//Check for each list what error was given from the target (if any):
		for _, ptn := range l {
			if ptn != "" && strings.Contains(body, ptn) {

				//Confirm that the error is unique and new to it's category:
				if !fc.InLst([]string{"\t", "\n", "\r", " "}, ptn) && !fc.InLst(vres_lFN, ptn) {
					lst = append(lst, ptn)
					//fmt.Println(count, "|", tkey, "->", lst) //DEBUG
					count++
				}
			}
		}

		//Items are collected from local gorutine. Adding to proper struct:
		result <- Edata{
			lErr:  lst,
			tErr:  tkey,
			count: count,
		}
	}
}
