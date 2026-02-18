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

package privategkemaster_contract

import (
	"testing"
)

func TestGKEMasterLogFieldSet_PrivateGKEMasterParserType(t *testing.T) {
	testCases := []struct {
		podID    string
		expected PrivateGKEMasterParserType
	}{
		{
			podID:    "kube-scheduler",
			expected: PrivateGKEMasterParserTypeScheduler,
		},
		{
			podID:    "kube-controller-manager",
			expected: PrivateGKEMasterParserTypeControllerManager,
		},
		{
			podID:    "kubelet",
			expected: PrivateGKEMasterParserTypeKubelet,
		},
		{
			podID:    "containerd",
			expected: PrivateGKEMasterParserTypeContainerRuntime,
		},
		{
			podID:    "kube-apiserver",
			expected: PrivateGKEMasterParserTypeOther,
		},
		{
			podID:    "random-component",
			expected: PrivateGKEMasterParserTypeOther,
		},
		{
			podID:    "",
			expected: PrivateGKEMasterParserTypeOther,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.podID, func(t *testing.T) {
			fs := &GKEMasterLogFieldSet{
				ComponentName: tc.podID,
			}
			got := fs.PrivateGKEMasterParserType()
			if got != tc.expected {
				t.Errorf("PrivateGKEMasterParserType() = %v, want %v", got, tc.expected)
			}
		})
	}
}
