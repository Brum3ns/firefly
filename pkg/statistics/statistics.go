// Provides general statistics for the runner process and all its nested processes/tasks.
package statistics

import (
	"time"
)

type Statistic struct {
	base
	Response       Http   `json:"Responses"`
	Request        Http   `json:"Request"`
	Output         Output `json:"Output"`
	Payload        `json:"Payload"`
	Scanner        `json:"Scanner"`
	Behavior       `json:"Behavior"`
	Pattern        `json:"Pattern"`
	Transformation `json:"Transformation"`
	Difference     `json:"Difference"`
}

type Http struct {
	base
	timeTotal   float64
	Timeout     int `json:"TimeTotal"`
	Forbidden   int `json:"Forbidden"`
	TimeAverage int `json:"TimeAverage"`
}

type Output struct{ base }
type Scanner struct{ base }
type Pattern struct{ base }
type Payload struct{ base }
type Behavior struct{ base }
type Difference struct{ base }
type Transformation struct{ base }

type base struct {
	err    int       `json:"error"`
	count  int       `json:"Count"`
	filter int       `json:"filter"`
	time   time.Time `json:"time"`
}

func (b *base) Count()              { b.count++ }
func (b *base) CountFilter()        { b.filter++ }
func (b *base) CountError()         { b.err++ }
func (b *base) GetErrorCount() int  { return b.err }
func (b *base) GetCount() int       { return b.count }
func (b *base) GetFilterCount() int { return b.filter }

func NewStatistic(verify bool) Statistic {
	return Statistic{
		base: base{time: time.Now()},
	}
}

// Set a timer
func (b *base) SetTime() time.Time {
	return time.Now()
}

// Return the time duration for how long the process has been running
func (b *base) GetTime() [3]time.Duration {
	t := time.Since(b.time)
	h := t / time.Hour
	t -= h * time.Hour
	m := t / time.Minute
	t -= m * time.Minute
	s := t / time.Second
	return [3]time.Duration{h, m, s}
}

func (h *Http) CountForbidden() {
	h.Forbidden++
}

func (h *Http) GetCountForbidden() int {
	return h.Forbidden
}

func (h *Http) UpdateTime(t float64) {
	h.timeTotal += t
}

func (h *Http) GetAverageTime() float64 {
	if h.timeTotal <= 0 || h.base.count <= 0 {
		return 0
	}
	return h.timeTotal / float64(h.base.count)
}
