// Provides general statistics for the runner process and all its nested processes/tasks.
package statistics

import (
	"time"
)

type Statistic struct {
	TotalRequests int  `json:"Jobs"` //!Note : (Must be declared all at once. Not one by one)
	Behavior      int  `json:"UnexpectedBehavior"`
	Error         int  `json:"Failed"`
	Completed     int  `json:"Completed"`
	Output        int  `json:"OutputCount"`
	Response      Data `json:"Responses"`
	Request       Data `json:"Request"`
	Scanner       Data `json:"Scanner"`
	Timer         time.Time
}

type Data struct {
	Count         int           `json:"Count"`
	Error         int           `json:"Error"`
	Filtered      int           `json:"Filtered"`
	ErrorMessages map[error]int `json:"ErrorMessages"`
}

type Skip struct {
	Tag string
	Err error
}

func NewStatistic(verify bool) Statistic {
	return Statistic{
		Timer: time.Now(),
	}
}

// Return the current status and a true boolean if it still have processes to handle, otherwise false.
func (st *Statistic) InProcess() bool {
	return !(st.TotalRequests > 0 && st.TotalRequests == st.GetTotal())
}

func (st *Statistic) ReqInSec() int {
	if st.TotalRequests <= 0 {
		return 0
	}
	return 0
}

func (st *Statistic) GetTotal() int {
	return (st.Completed + st.Error + st.Request.Filtered)
}

func (st *Statistic) CountOutput() {
	st.Output++
}

func (st *Statistic) CountComplete() {
	st.Completed++
}

func (st *Statistic) CountError() {
	st.Error++
}

func (st *Statistic) CountBehavior() {
	st.Behavior++
}

func (st *Statistic) CountFilter() {
	st.Request.Filtered++
}

func (st *Statistic) CountRequest() {
	st.Request.Count++
}

func (st *Statistic) CountResponse() {
	st.Response.Count++
}

func (st *Statistic) CountScanner() {
	st.Scanner.Count++
}

// Return the time on how long the process has been running
func (st *Statistic) getTime() [3]time.Duration {
	t := time.Since(st.Timer)
	h := t / time.Hour
	t -= h * time.Hour
	m := t / time.Minute
	t -= m * time.Minute
	s := t / time.Second
	return [3]time.Duration{h, m, s}
}
