package common

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentCounter(t *testing.T) {
	ThreadCount := 100
	ItemCount := 100
	counter := NewDefaultConcurrentCounter(NewSuffixShardingProvider(16, 1))
	wg := sync.WaitGroup{}
	for tc := 0; tc < ThreadCount; tc++ {
		wg.Add(1)
		go func() {
			for ic := 0; ic < ItemCount; ic++ {
				counter.Incr(fmt.Sprintf("item-%d", ic))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	for ic := 0; ic < ItemCount; ic++ {
		cnt := counter.Get(fmt.Sprintf("item-%d", ic))
		if cnt != ItemCount {
			t.Errorf("expected %d, got %d", ItemCount, cnt)
		}
	}
}
