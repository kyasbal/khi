package progress

import (
	"context"
	"fmt"
	"time"
)

type ProgressUpdatorOnTickFunc = func(tp *TaskProgress)

type ProgressUpdator struct {
	Progress *TaskProgress
	Interval time.Duration
	OnTick   ProgressUpdatorOnTickFunc
	context  context.Context
	cancel   func()
}

func NewProgressUpdator(progress *TaskProgress, interval time.Duration, onTick ProgressUpdatorOnTickFunc) *ProgressUpdator {
	return &ProgressUpdator{
		Progress: progress,
		Interval: interval,
		OnTick:   onTick,
	}
}

func (p *ProgressUpdator) Start(ctx context.Context) error {
	p.OnTick(p.Progress)
	cancellable, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.context = cancellable
	go func() {
		for itr := 1; true; itr++ {
			select {
			case <-p.context.Done():
				return
			case <-time.After(p.Interval):
				p.OnTick(p.Progress)
				itr++
			}
		}
	}()
	return nil
}

func (p *ProgressUpdator) Done() error {
	if p.context == nil {
		return fmt.Errorf("this updator is not yet started")
	}
	p.cancel()
	return nil
}
