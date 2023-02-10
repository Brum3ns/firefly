package runner

import (
	"sync"

	"github.com/Brum3ns/firefly/pkg/firefly/technique"
	"github.com/Brum3ns/firefly/pkg/functions"
	G "github.com/Brum3ns/firefly/pkg/functions/globalVariables"
	"github.com/Brum3ns/firefly/pkg/storage"
)

/* func (r *Runner) Calcs(ID int, resp storage.Response) *EiCalc {
	var (
		lc = functions.LineIn(resp.Body)
		wc = functions.WordIn(resp.Body)
	)
	return &EiCalc{
		LineCount: lc,
		WordCount: wc,

		ID:     ID,
		Status: true,
	}
} */

func (r *Runner) Transformation(ID int, resp storage.Response) *EiTransformation {
	/** Analyze the target for payload transformations */
	var (
		ok bool
		m  = make(map[string]string)
	)

	if !r.verify && resp.Tag == "transformation" {
		m, ok = technique.TransformationScan(r.wordlist.TransformationCompare, resp)
	} else {
		m = nil
		ok = false
	}

	return &EiTransformation{
		Transformation:         m,
		Transformation_display: functions.MapToString(m, " → "),
		OK:                     ok,
		ID:                     ID,
		Status:                 true,
	}
}

func (r *Runner) Difference(ID int, resp storage.Response) *EiDifference {
	/** Behavior detection from all response sources
	* Return data that contains behaviors related to the target.
	* If a payload effected the target behavior it analyzed it.
	 */

	var (

		//[TODO] Move to diff:
		lc = functions.LineIn(resp.Body)
		wc = functions.WordIn(resp.Body)

		//amount int
		//m           = make(map[string]map[string]string)
		m_AvgAmount = make(map[string]int)
		//m_tmp  = make(map[string]string)
	)

	if r.verify {
		//m = nil
	} else {
		_, m_AvgAmount = technique.Diff("body", resp, r.verifyData)
		//m, banner = technique.Diff("body", resp, r.verifyData) //<-- TEMP OLD
	}

	return &EiDifference{
		ID: ID,
		//Diff:      m,
		AvgAmount: m_AvgAmount,
		WordCount: wc,
		LineCount: lc,
		Status:    true,
	}

}

func (r *Runner) Error(ID int, resp storage.Response) *EiErrors {
	/** Extract errors from response(s)
	* Return collected errors & sort it into error category
	 */

	m, lstErrors, hasErr, ok := technique.ErrorDetect(string(resp.Body[:]), r.wordlist, r.verifyData)

	return &EiErrors{
		ID: ID,
		//OK:     hit,
		Errs:       m,
		OK:         ok,
		ErrsAmount: lstErrors,
		HasErr:     hasErr,
		Status:     true,
	}
}

//[TODO] In development...
func (r *Runner) Pattern(ID int) *EiPatterns {
	/** Detect & collect Patterns from responses
	* Return extracted patterns & sort each pattern
	* into it's category.
	 */
	//m := make(chan *storage.Wordlists)

	return &EiPatterns{
		ID:     ID,
		Status: true,
	}
}

/**Filter - Match 'x' type. Return 'True' if it match the given data*/
func (r *Runner) Filter(rt storage.Response) bool {
	/**More can be found inside runner 'filter.go'*/
	var (
		mutex    = &sync.Mutex{}
		okF      bool
		okM      bool
		Amount   = 0
		m_Filter = make(map[string]map[string]bool)
	)

	for _, mode := range []string{"filter", "match"} {
		var m_Mode map[string][]string

		switch mode {
		case "filter":
			m_Mode = G.Lst_mFilter
		case "match":
			m_Mode = G.Lst_mMatch
		}
		total := len(m_Mode)

		//Do not continue if 'x' mode isn't set:
		if total <= 0 {
			continue
		}

		//setup [chan]nels:
		c_jobs := make(chan FData, total)
		c_result := make(chan FResult, total)
		defer close(c_jobs)
		defer close(c_result)

		//Start 'filter' task & Give the jobs: (filter that should be used)
		for t := 0; t < total; t++ {
			go FilterCheck(rt, c_jobs, c_result)
		}
		JobsFilter(m_Mode, c_jobs)

		//Check result
		for i := 0; i < total; i++ {
			result := <-c_result

			mutex.Lock()
			m_Filter[mode] = make(map[string]bool)
			m_Filter[mode] = result.m

			if m_Filter[mode][result.t] {
				Amount++
			}
			mutex.Unlock()
		}

		//Check if all filter(s) factors are true:
		switch mode {
		case "filter":
			if Amount == total {
				okF = true
			}
		case "match":
			if Amount == total {
				okM = true
			}
		}
	}

	// True='skip', False='continue': [TODO] - (Improve code structure)
	var (
		M  = len(m_Filter["match"])
		F  = len(m_Filter["filter"])
		ok bool
	)
	if F > 0 && okF {
		ok = true
	} else if (F > 0 && !okF) && (M <= 0 && !okM) {
		ok = false
	} else if (F <= 0 || !okF) && (M >= 0 && okM) {
		ok = false
	} else if (F <= 0 || !okF) && (M >= 0 && !okM) {
		ok = true
	}

	return ok
}
