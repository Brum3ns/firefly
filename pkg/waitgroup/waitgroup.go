package waitgroup

import (
	"sync"
	"sync/atomic"
)

type WaitGroup struct {
	sync.WaitGroup
	count int64
}

func (wg *WaitGroup) Add(delta int) {
	atomic.AddInt64(&wg.count, int64(delta))
	wg.WaitGroup.Add(delta)
}

func (wg *WaitGroup) Done() {
	atomic.AddInt64(&wg.count, -1)
	wg.WaitGroup.Done()
}

func (wg *WaitGroup) GetCount() int {
	return int(atomic.LoadInt64(&wg.count))
}
