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

package server

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestResourceMonitorImpl(t *testing.T) {
	monitor := &ResourceMonitorImpl{}

	testCases := []struct {
		name string
		op   func() int
	}{
		{
			name: "GetUsedMemory",
			op:   monitor.GetUsedMemory,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.op()
			if got <= 0 {
				t.Errorf("%s() = %d; want > 0", tc.name, got)
			}
		})
	}
}

func TestResourceMonitorMock(t *testing.T) {
	testCases := []struct {
		name     string
		mock     ResourceMonitorMock
		wantUsed int
	}{
		{
			name: "mock returns set values",
			mock: ResourceMonitorMock{
				UsedMemory: 100,
			},
			wantUsed: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if diff := cmp.Diff(tc.wantUsed, tc.mock.GetUsedMemory()); diff != "" {
				t.Errorf("GetUsedMemory() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
