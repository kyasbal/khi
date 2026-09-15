// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inspectionmetadata

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetTaskProgress(t *testing.T) {
	progress := NewProgress()
	tp, err := progress.GetOrCreateTaskProgress("foo")
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	expected := TaskProgressSnapshot{
		ID:      "foo",
		Ratio:   0,
		Message: "",
		Label:   "foo",
	}

	if diff := cmp.Diff(expected, tp.Snapshot()); diff != "" {
		t.Errorf("generated task progress is not containing the expected state\n%s", diff)
	}

	tp2, err := progress.GetOrCreateTaskProgress("foo")
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	if tp != tp2 {
		t.Errorf("GetTaskProgress should return the same reference for the existing task")
	}
}

func TestTaskProgressMetadataMethods(t *testing.T) {
	testCases := []struct {
		name   string
		action func(tp *TaskProgressMetadata)
		want   TaskProgressSnapshot
	}{
		{
			name: "set label and ratio update",
			action: func(tp *TaskProgressMetadata) {
				tp.SetLabel("Custom Label")
				tp.Update(0.45, "Processing items")
			},
			want: TaskProgressSnapshot{
				ID:            "task-1",
				Label:         "Custom Label",
				Message:       "Processing items",
				Ratio:         0.45,
				Indeterminate: false,
			},
		},
		{
			name: "update indeterminate resets ratio",
			action: func(tp *TaskProgressMetadata) {
				tp.SetLabel("Custom Label")
				tp.Update(0.45, "Processing items")
				tp.UpdateIndeterminate("Searching...")
			},
			want: TaskProgressSnapshot{
				ID:            "task-1",
				Label:         "Custom Label",
				Message:       "Searching...",
				Ratio:         0,
				Indeterminate: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tp := NewTaskProgressMetadata("task-1")
			tc.action(tp)
			if diff := cmp.Diff(tc.want, tp.Snapshot()); diff != "" {
				t.Errorf("TaskProgressMetadata.Snapshot() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProgressPhaseTransitions(t *testing.T) {
	testCases := []struct {
		name            string
		totalTasks      int
		action          func(p *Progress) error
		want            ProgressSnapshot
		isTerminalPhase bool
	}{
		{
			name:       "resolve single task",
			totalTasks: 2,
			action: func(p *Progress) error {
				if _, err := p.GetOrCreateTaskProgress("foo"); err != nil {
					return err
				}
				if _, err := p.GetOrCreateTaskProgress("bar"); err != nil {
					return err
				}
				return p.ResolveTask("foo")
			},
			want: ProgressSnapshot{
				Phase:         TaskPhaseRunning,
				TotalProgress: &TaskProgressSnapshot{ID: "Total", Label: "Total", Message: "1 of 2 tasks complete", Ratio: 0.5},
				TaskProgresses: []TaskProgressSnapshot{
					{ID: "bar", Label: "bar"},
				},
			},
			isTerminalPhase: false,
		},
		{
			name:       "mark done clears tasks and blocks subsequent mutations",
			totalTasks: 2,
			action: func(p *Progress) error {
				if _, err := p.GetOrCreateTaskProgress("foo"); err != nil {
					return err
				}
				if _, err := p.GetOrCreateTaskProgress("bar"); err != nil {
					return err
				}
				return p.MarkDone()
			},
			want: ProgressSnapshot{
				Phase:          TaskPhaseDone,
				TotalProgress:  &TaskProgressSnapshot{ID: "Total", Label: "Total", Message: "2 of 2 tasks complete", Ratio: 1},
				TaskProgresses: []TaskProgressSnapshot{},
			},
			isTerminalPhase: true,
		},
		{
			name:       "mark done with zero tasks",
			totalTasks: 0,
			action: func(p *Progress) error {
				return p.MarkDone()
			},
			want: ProgressSnapshot{
				Phase:          TaskPhaseDone,
				TotalProgress:  &TaskProgressSnapshot{ID: "Total", Label: "Total", Message: "Complete", Ratio: 1},
				TaskProgresses: []TaskProgressSnapshot{},
			},
			isTerminalPhase: true,
		},
		{
			name:       "mark cancelled clears tasks and blocks subsequent mutations",
			totalTasks: 2,
			action: func(p *Progress) error {
				if _, err := p.GetOrCreateTaskProgress("foo"); err != nil {
					return err
				}
				if _, err := p.GetOrCreateTaskProgress("bar"); err != nil {
					return err
				}
				return p.MarkCancelled()
			},
			want: ProgressSnapshot{
				Phase:          TaskPhaseCancelled,
				TotalProgress:  &TaskProgressSnapshot{ID: "Total", Label: "Total", Message: "0 of 2 tasks complete"},
				TaskProgresses: []TaskProgressSnapshot{},
			},
			isTerminalPhase: true,
		},
		{
			name:       "mark error clears tasks and blocks subsequent mutations",
			totalTasks: 2,
			action: func(p *Progress) error {
				if _, err := p.GetOrCreateTaskProgress("foo"); err != nil {
					return err
				}
				return p.MarkError()
			},
			want: ProgressSnapshot{
				Phase:          TaskPhaseError,
				TotalProgress:  &TaskProgressSnapshot{ID: "Total", Label: "Total", Message: "0 of 2 tasks complete"},
				TaskProgresses: []TaskProgressSnapshot{},
			},
			isTerminalPhase: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewProgress()
			if tc.totalTasks > 0 {
				p.SetTotalTaskCount(tc.totalTasks)
			}
			if err := tc.action(p); err != nil {
				t.Fatalf("action failed: %v", err)
			}
			if diff := cmp.Diff(tc.want, p.Snapshot()); diff != "" {
				t.Errorf("Progress.Snapshot() mismatch (-want +got):\n%s", diff)
			}
			if tc.isTerminalPhase {
				if _, err := p.GetOrCreateTaskProgress("after"); err == nil {
					t.Errorf("GetOrCreateTaskProgress() after terminal phase = nil, want error")
				}
				if err := p.ResolveTask("after"); err == nil {
					t.Errorf("ResolveTask() after terminal phase = nil, want error")
				}
				if err := p.MarkDone(); err == nil {
					t.Errorf("MarkDone() after terminal phase = nil, want error")
				}
				if err := p.MarkCancelled(); err == nil {
					t.Errorf("MarkCancelled() after terminal phase = nil, want error")
				}
				if err := p.MarkError(); err == nil {
					t.Errorf("MarkError() after terminal phase = nil, want error")
				}
			}
		})
	}
}

func TestProgressSnapshotAndMarshalJSON(t *testing.T) {
	progress := NewProgress()
	progress.SetTotalTaskCount(2)
	tp, err := progress.GetOrCreateTaskProgress("task-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tp.SetLabel("Task A")
	tp.Update(0.5, "Working")

	snap := progress.Snapshot()
	wantSnap := ProgressSnapshot{
		Phase: TaskPhaseRunning,
		TotalProgress: &TaskProgressSnapshot{
			ID:      "Total",
			Label:   "Total",
			Message: "0 of 2 tasks complete",
		},
		TaskProgresses: []TaskProgressSnapshot{
			{
				ID:      "task-a",
				Label:   "Task A",
				Message: "Working",
				Ratio:   0.5,
			},
		},
	}
	if diff := cmp.Diff(wantSnap, snap); diff != "" {
		t.Errorf("Snapshot() mismatch (-want +got):\n%s", diff)
	}

	rawJSON, err := json.Marshal(progress)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	var decoded ProgressSnapshot
	if err := json.Unmarshal(rawJSON, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if diff := cmp.Diff(wantSnap, decoded); diff != "" {
		t.Errorf("MarshalJSON() roundtrip mismatch (-want +got):\n%s", diff)
	}
}
