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

package private

import (
	"context"
	"testing"

	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	"github.com/google/go-cmp/cmp"
)

func TestPrivateLifecycleInitializer(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "registers OnRun hook for lifecycle notification",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			engine := coreinit.NewEngine(context.Background())
			ctx := engine.Context()

			if err := PrivateLifecycleInitializer.Init(ctx); err != nil {
				t.Fatalf("PrivateLifecycleInitializer.Init() failed: %v", err)
			}

			if diff := cmp.Diff(PrivateLifecycleInitializer.ID, InitializerIDPrivateLifecycle); diff != "" {
				t.Errorf("Initializer ID mismatch (-want +got):\n%s", diff)
			}

			if err := engine.Run(); err != nil {
				t.Fatalf("Engine.Run() failed: %v", err)
			}
		})
	}
}
