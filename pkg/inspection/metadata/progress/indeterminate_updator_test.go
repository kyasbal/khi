package progress

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestIndeterminateUpdator(t *testing.T) {
	progress := NewTaskProgress("foo")
	updator := NewIndeterminateUpdator(progress, 1000*time.Millisecond)
	err := updator.Start("working")
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	time.Sleep(1500 * time.Millisecond)
	if diff := cmp.Diff(&TaskProgress{
		Id:            "foo",
		Label:         "foo",
		Message:       "working.",
		Percentage:    0,
		Indeterminate: true,
	}, progress, cmpopts.IgnoreUnexported(TaskProgress{})); diff != "" {
		t.Errorf("The result status is not in the expected status\n%s", diff)
	}
	err = updator.Done()
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	msg := progress.Message
	time.Sleep(1000 * time.Millisecond)
	if msg != progress.Message {
		t.Errorf("The progress is keeping updated")
	}
}
