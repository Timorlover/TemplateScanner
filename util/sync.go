package util

import "sync"

type WaitGroup struct {
	size      int
	pool      chan byte
	waitCount int64
	waitGroup sync.WaitGroup
}

func NewWaitGroup(size int) *WaitGroup {
	wg := &WaitGroup{
		size: size,
	}
	if size > 0 {
		wg.pool = make(chan byte, size)
	}
	return wg
}

func (wg *WaitGroup) BlockAdd() {
	if wg.size > 0 {
		wg.pool <- 1
	}
	wg.waitGroup.Add(1)
}

func (wg *WaitGroup) Done() {
	if wg.size > 0 {
		<-wg.pool
	}
	wg.waitGroup.Done()
}

func (wg *WaitGroup) Wait() {
	wg.waitGroup.Wait()
}

func (wg *WaitGroup) PendingCount() int64 {
	return int64(len(wg.pool))
}
