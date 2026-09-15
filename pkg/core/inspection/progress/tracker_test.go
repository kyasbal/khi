// Copyright 2026 Google LLC
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

package progress

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
)

func TestFormatETA(t *testing.T) {
	testCases := []struct {
		name  string
		input time.Duration
		want  string
	}{
		{name: "negative duration", input: -5 * time.Second, want: "0s"},
		{name: "zero duration", input: 0, want: "0s"},
		{name: "seconds under one minute", input: 45 * time.Second, want: "45s"},
		{name: "rounding seconds", input: 14*time.Second + 600*time.Millisecond, want: "15s"},
		{name: "exact one minute", input: 60 * time.Second, want: "1m00s"},
		{name: "minutes and single-digit seconds", input: 65 * time.Second, want: "1m05s"},
		{name: "minutes and double-digit seconds", input: 95 * time.Second, want: "1m35s"},
		{name: "exact one hour", input: 3600 * time.Second, want: "1h00m"},
		{name: "hours and single-digit minutes", input: 3600*time.Second + 5*time.Minute, want: "1h05m"},
		{name: "hours and double-digit minutes", input: 2*time.Hour + 30*time.Minute, want: "2h30m"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := formatETA(tc.input)
			if got != tc.want {
				t.Errorf("formatETA(%v) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestCalculateETA(t *testing.T) {
	testCases := []struct {
		name          string
		elapsed       time.Duration
		completeRatio float32
		wantETA       string
	}{
		{
			name:          "elapsed under 1 second returns placeholder",
			elapsed:       500 * time.Millisecond,
			completeRatio: 0.5,
			wantETA:       "--",
		},
		{
			name:          "ratio zero returns placeholder",
			elapsed:       10 * time.Second,
			completeRatio: 0,
			wantETA:       "--",
		},
		{
			name:          "ratio under 0.005 returns placeholder",
			elapsed:       10 * time.Second,
			completeRatio: 0.001,
			wantETA:       "--",
		},
		{
			name:          "100 percent complete returns 0s",
			elapsed:       10 * time.Second,
			completeRatio: 1.0,
			wantETA:       "0s",
		},
		{
			name:          "ratio exceeding 1.0 returns 0s",
			elapsed:       10 * time.Second,
			completeRatio: 1.2,
			wantETA:       "0s",
		},
		{
			name:          "50 percent complete after 10s returns 10s",
			elapsed:       10 * time.Second,
			completeRatio: 0.5,
			wantETA:       "10s",
		},
		{
			name:          "25 percent complete after 30s returns 1m30s",
			elapsed:       30 * time.Second,
			completeRatio: 0.25,
			wantETA:       "1m30s",
		},
		{
			name:          "10 percent complete after 10 minutes returns 1h30m",
			elapsed:       10 * time.Minute,
			completeRatio: 0.1,
			wantETA:       "1h30m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculateETA(tc.elapsed, tc.completeRatio)
			if got != tc.wantETA {
				t.Errorf("calculateETA(%v, %v) = %q, want %q", tc.elapsed, tc.completeRatio, got, tc.wantETA)
			}
		})
	}
}

func TestTrackerLifecycle(t *testing.T) {
	testCases := []struct {
		name              string
		total             int
		opts              []TrackerOption
		actions           func(tr *Tracker)
		wantRatio         float32
		wantIndeterminate bool
		wantMsgSubstring  string
	}{
		{
			name:  "discrete progress with unit",
			total: 10,
			opts: []TrackerOption{
				WithUnit("logs"),
			},
			actions: func(tr *Tracker) {
				tr.Add(4)
				tr.Inc()
				tr.tick()
			},
			wantRatio:         0.5,
			wantIndeterminate: false,
			wantMsgSubstring:  "5/10 logs",
		},
		{
			name:  "elapsed over 1s formats throughput rate and ETA",
			total: 10,
			opts: []TrackerOption{
				WithUnit("logs"),
			},
			actions: func(tr *Tracker) {
				if tr.cancel != nil {
					tr.cancel()
					tr.wg.Wait()
				}
				tr.startTime = time.Now().Add(-2 * time.Second)
				tr.Add(5)
				tr.tick()
			},
			wantRatio:         0.5,
			wantIndeterminate: false,
			wantMsgSubstring:  "lps, ETA",
		},
		{
			name:  "tracker done flushes completed count",
			total: 5,
			opts: []TrackerOption{
				WithUnit("items"),
			},
			actions: func(tr *Tracker) {
				tr.Add(5)
				tr.Done()
			},
			wantRatio:         1.0,
			wantIndeterminate: false,
			wantMsgSubstring:  "5/5 items",
		},
		{
			name:  "clamps count exceeding total and handles multiple Done calls",
			total: 5,
			opts: []TrackerOption{
				WithUnit("items"),
			},
			actions: func(tr *Tracker) {
				tr.Add(8)
				tr.Done()
				tr.Done()
			},
			wantRatio:         1.0,
			wantIndeterminate: false,
			wantMsgSubstring:  "5/5 items",
		},
		{
			name:              "zero total starts indeterminate before done is called",
			total:             0,
			opts:              nil,
			actions:           func(tr *Tracker) {},
			wantRatio:         0,
			wantIndeterminate: true,
			wantMsgSubstring:  "",
		},
		{
			name:  "zero total starts indeterminate and resolves on done",
			total: 0,
			opts:  nil,
			actions: func(tr *Tracker) {
				tr.Done()
			},
			wantRatio:         1.0,
			wantIndeterminate: false,
			wantMsgSubstring:  "Complete",
		},
		{
			name:  "discrete progress without unit option formats raw count and items/s",
			total: 10,
			opts:  nil,
			actions: func(tr *Tracker) {
				if tr.cancel != nil {
					tr.cancel()
					tr.wg.Wait()
				}
				tr.startTime = time.Now().Add(-2 * time.Second)
				tr.Add(5)
				tr.tick()
			},
			wantRatio:         0.5,
			wantIndeterminate: false,
			wantMsgSubstring:  "5/10 (2.5 items/s, ETA",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tp := inspectionmetadata.NewTaskProgressMetadata("task-id")
			ctx := WithContext(context.Background(), tp)

			tr := NewTracker(ctx, tc.total, tc.opts...)
			defer tr.Done()
			tc.actions(tr)

			snap := tp.Snapshot()
			if snap.Ratio != tc.wantRatio {
				t.Errorf("snap.Ratio = %v, want %v", snap.Ratio, tc.wantRatio)
			}
			if snap.Indeterminate != tc.wantIndeterminate {
				t.Errorf("snap.Indeterminate = %v, want %v", snap.Indeterminate, tc.wantIndeterminate)
			}
			if !strings.Contains(snap.Message, tc.wantMsgSubstring) {
				t.Errorf("snap.Message = %q, want substring %q", snap.Message, tc.wantMsgSubstring)
			}
		})
	}
}

func TestForEach(t *testing.T) {
	testCases := []struct {
		name          string
		items         []int
		cancelAtIndex int
		failAtIndex   int
		wantErr       bool
		wantRatio     float32
	}{
		{
			name:          "completes all items successfully",
			items:         []int{10, 20, 30, 40},
			cancelAtIndex: -1,
			failAtIndex:   -1,
			wantErr:       false,
			wantRatio:     1.0,
		},
		{
			name:          "empty slice completes immediately with ratio 1.0",
			items:         []int{},
			cancelAtIndex: -1,
			failAtIndex:   -1,
			wantErr:       false,
			wantRatio:     1.0,
		},
		{
			name:          "stops early on error and preserves actual progress",
			items:         []int{10, 20, 30, 40},
			cancelAtIndex: -1,
			failAtIndex:   2,
			wantErr:       true,
			wantRatio:     0.5,
		},
		{
			name:          "stops early when context is cancelled",
			items:         []int{10, 20, 30, 40},
			cancelAtIndex: 1,
			failAtIndex:   -1,
			wantErr:       true,
			wantRatio:     0.5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tp := inspectionmetadata.NewTaskProgressMetadata("foreach-task")
			ctx, cancel := context.WithCancel(WithContext(context.Background(), tp))
			defer cancel()

			err := ForEach(ctx, tc.items, func(index int, item int) error {
				if index == tc.cancelAtIndex {
					cancel()
				}
				if index == tc.failAtIndex {
					return errors.New("processing error")
				}
				return nil
			}, WithUnit("lines"))

			if (err != nil) != tc.wantErr {
				t.Errorf("ForEach() error = %v, wantErr %v", err, tc.wantErr)
			}
			snap := tp.Snapshot()
			if snap.Ratio != tc.wantRatio {
				t.Errorf("snap.Ratio = %v, want %v", snap.Ratio, tc.wantRatio)
			}
		})
	}
}

func TestRatioTracker(t *testing.T) {
	testCases := []struct {
		name             string
		opts             []TrackerOption
		backdateElapsed  time.Duration
		skipUpdate       bool
		ratio            float32
		count            int
		updateOpts       []RatioUpdateOption
		wantRatio        float32
		wantMsgSubstring string
	}{
		{
			name:             "sub-second update with step and log unit",
			opts:             []TrackerOption{WithUnit("logs")},
			ratio:            0.25,
			count:            500,
			updateOpts:       []RatioUpdateOption{WithStep(1, 4)},
			wantRatio:        0.25,
			wantMsgSubstring: "[1/4] 500 logs",
		},
		{
			name:             "elapsed over 1s with non-log unit formats chunks/s and ETA",
			opts:             []TrackerOption{WithUnit("chunks")},
			backdateElapsed:  2 * time.Second,
			ratio:            0.5,
			count:            20,
			wantRatio:        0.5,
			wantMsgSubstring: "chunks/s, ETA",
		},
		{
			name:             "elapsed over 1s with empty unit formats items/s",
			opts:             nil,
			backdateElapsed:  2 * time.Second,
			ratio:            1.0,
			count:            100,
			wantRatio:        1.0,
			wantMsgSubstring: "items/s, ETA 0s",
		},
		{
			name:             "zero updates before Done marks Complete",
			opts:             []TrackerOption{WithUnit("logs")},
			skipUpdate:       true,
			wantRatio:        1.0,
			wantMsgSubstring: "Complete",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tp := inspectionmetadata.NewTaskProgressMetadata("ratio-task")
			ctx := WithContext(context.Background(), tp)

			rt := NewRatioTracker(ctx, tc.opts...)
			if tc.backdateElapsed > 0 {
				rt.startTime = time.Now().Add(-tc.backdateElapsed)
			}
			if !tc.skipUpdate {
				rt.Update(tc.ratio, tc.count, tc.updateOpts...)

				snap := tp.Snapshot()
				if snap.Ratio != tc.wantRatio {
					t.Errorf("snap.Ratio = %v, want %v", snap.Ratio, tc.wantRatio)
				}
				if !strings.Contains(snap.Message, tc.wantMsgSubstring) {
					t.Errorf("snap.Message = %q, want substring %q", snap.Message, tc.wantMsgSubstring)
				}
			}

			rt.Done()
			snap := tp.Snapshot()
			if snap.Ratio != tc.wantRatio {
				t.Errorf("after Done(), snap.Ratio = %v, want %v", snap.Ratio, tc.wantRatio)
			}
			if !strings.Contains(snap.Message, tc.wantMsgSubstring) {
				t.Errorf("after Done(), snap.Message = %q, want substring %q", snap.Message, tc.wantMsgSubstring)
			}
		})
	}
}
