// Provides general statistics for the runner process and all its nested processes/tasks.
package runner

import (
	"sync"
)

type Statistic struct {
	TotalRequests      int   `json:"Jobs"` //!Note : (Must be declared all at once. Not one by one)
	UnexpectedBehavior int   `json:"UnexpectedBehavior"`
	Failed             int   `json:"Failed"`
	Response           int   `json:"Responses"`
	Completed          int   `json:"Completed"`
	Output             int   `json:"OutputCount"`
	Request            *Data `json:"Request"`
	Scanner            *Data `json:"Scanner"`
	Mutex              sync.Mutex
}
type Data struct {
	Count         int           `json:"Count"`
	Error         int           `json:"Error"`
	Filtered      int           `json:"Filtered"`
	ErrorMessages map[error]int `json:"ErrorMessages"`
}

func newStatistic() *Statistic {
	return &Statistic{
		Request: &Data{ErrorMessages: make(map[error]int)},
		Scanner: &Data{ErrorMessages: make(map[error]int)},
	}
}

func (st *Statistic) ReqInSec() int {
	if st.TotalRequests <= 0 {
		return 0
	}
	return 0
}

// Return the current status and a true boolean if it still have processes to handle, otherwise false.
func (st *Statistic) inProcess() bool {
	return !(st.TotalRequests > 0 && st.TotalRequests == st.getTotal())
}

func (st *Statistic) handleSkipped(skip skipProcess) {
	if skip.tag == "filter" {
		st.countFilter()
	} else if skip.tag == "error" {
		st.countFail()
		if skip.err != nil {
			st.appendRequestError(skip.err)
		}
	}
}

func (st *Statistic) getTotal() int {
	return (st.Completed + st.Failed + st.Request.Filtered)
}

func (st *Statistic) countOutput() {
	st.Output++
}

func (st *Statistic) countComplete() {
	st.Completed++
}

func (st *Statistic) countFail() {
	st.Failed++
}

func (st *Statistic) countUnexpectedBehavior() {
	st.UnexpectedBehavior++
}

func (st *Statistic) countFilter() {
	st.Request.Filtered++
}

func (st *Statistic) countRequest() {
	st.Request.Count++
}

func (st *Statistic) countResponse() {
	st.Response++
}

func (st *Statistic) countScanner() {
	st.Scanner.Count++
}

// Add a new request error to the runner statistics instant
func (st *Statistic) appendRequestError(err error) {
	if err != nil {
		st.Mutex.Lock()
		st.Request.ErrorMessages[err]++
		st.Mutex.Unlock()
		st.Request.Error++
	}
}

// Add a new scanner error to the runner statistics instant
func (st *Statistic) appendScannerError(err error) {
	if err != nil {
		st.Mutex.Lock()
		st.Scanner.ErrorMessages[err]++
		st.Mutex.Unlock()
		st.Scanner.Error++
	}
}
