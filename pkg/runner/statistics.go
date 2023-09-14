// Provides general statistics for the runner process and all its nested processes/tasks.
package runner

import (
	"sync"
)

type Statistics struct {
	jobTotal  int   `json:"Jobs"` //!Note : (Must be declared all at once. Not one by one)
	Failed    int   `json:"Failed"`
	Responses int   `json:"Responses"`
	Completed int   `json:"Completed"`
	Output    int   `json:"OutputCount"`
	Request   *Data `json:"Request"`
	Scanner   *Data `json:"Scanner"`
	Mutex     sync.Mutex
}
type Data struct {
	Count         int           `json:"Count"`
	Error         int           `json:"Error"`
	Filtered      int           `json:"Filtered"`
	ErrorMessages map[error]int `json:"ErrorMessages"`
}

func newStatistic() *Statistics {
	return &Statistics{
		Request: &Data{ErrorMessages: make(map[error]int)},
		Scanner: &Data{ErrorMessages: make(map[error]int)},
	}
}

// Return the current status and a true boolean if it still have processes to handle, otherwise false.
func (st *Statistics) inProcess() bool {
	return !(st.jobTotal > 0 && st.jobTotal == st.getTotal())
}

func (st *Statistics) handleSkipped(skip skipProcess) {
	if skip.tag == "filter" {
		st.countFilter()
	} else if skip.tag == "error" {
		st.countFail()
		if skip.err != nil {
			st.appendRequestError(skip.err)
		}
	}
}

func (st *Statistics) getTotal() int {
	return (st.Completed + st.Failed + st.Request.Filtered)
}

func (st *Statistics) countOutput() {
	st.Output++
}

func (st *Statistics) countComplete() {
	st.Completed++
}

func (st *Statistics) countFail() {
	st.Failed++
}

func (st *Statistics) countFilter() {
	st.Request.Filtered++
}

func (st *Statistics) countRequest() {
	st.Request.Count++
}

func (st *Statistics) countScanner() {
	st.Scanner.Count++
}

// Add a new request error to the runner statistics instant
func (st *Statistics) appendRequestError(err error) {
	if err != nil {
		st.Mutex.Lock()
		st.Request.ErrorMessages[err]++
		st.Mutex.Unlock()
		st.Request.Error++
	}
}

// Add a new scanner error to the runner statistics instant
func (st *Statistics) appendScannerError(err error) {
	if err != nil {
		st.Mutex.Lock()
		st.Scanner.ErrorMessages[err]++
		st.Mutex.Unlock()
		st.Scanner.Error++
	}
}
