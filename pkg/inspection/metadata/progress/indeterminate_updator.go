package progress

import (
	"context"
	"fmt"
	"time"
)

// IndeterminateUpdator updates progress bar during a procedure can't report its progress.
type IndeterminateUpdator struct {
	Progress *TaskProgress
	Interval time.Duration
	context  context.Context
	cancel   func()
}

func NewIndeterminateUpdator(progress *TaskProgress, interval time.Duration) *IndeterminateUpdator {
	progress.Indeterminate = true
	return &IndeterminateUpdator{
		Progress: progress,
		Interval: interval,
	}
}

// Start starts updating progress bar.
// It returns an error if the updator is already started.
func (i *IndeterminateUpdator) Start(msg string) error {
	if i.context != nil {
		return fmt.Errorf("this updator is already used")
	}
	cancellable, cancel := context.WithCancel(context.Background())
	i.Progress.Message = msg
	i.context = cancellable
	i.cancel = cancel
	go func() {
		for itr := 1; true; itr++ {
			select {
			case <-i.context.Done():
				i.Progress.Indeterminate = false
				return
			case <-time.After(i.Interval):
				i.Progress.Message = fmt.Sprintf("%s%s", msg, i.workingIndicator(itr))
				itr++
			}
		}
	}()
	return nil
}

// Done stops updating progress bar.
// It returns an error if the updator is not yet started.
func (i *IndeterminateUpdator) Done() error {
	if i.context == nil {
		return fmt.Errorf("this updator is not yet started")
	}
	i.cancel()
	return nil
}

func (i *IndeterminateUpdator) workingIndicator(itr int) string {
	dots := ""
	for i := 0; i < itr%20; i++ {
		dots += "."
	}
	return dots
}
