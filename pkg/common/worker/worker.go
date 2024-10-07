package worker

import "sync"

// Pool enables running a goroutine with max parallel count limit.
type Pool struct {
	semaphore chan struct{}
	waitGroup *sync.WaitGroup
}

func NewPool(maxParallelCount int) *Pool {
	return &Pool{
		semaphore: make(chan struct{}, maxParallelCount),
		waitGroup: &sync.WaitGroup{},
	}
}

func (t *Pool) Run(f func()) {
	t.waitGroup.Add(1)
	t.semaphore <- struct{}{}
	go func() {
		defer func() {
			<-t.semaphore
			t.waitGroup.Done()
		}()
		f()
	}()
}

func (t *Pool) Wait() {
	t.waitGroup.Wait()
}
