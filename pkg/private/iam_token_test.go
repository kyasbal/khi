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

	"github.com/GoogleCloudPlatform/khi/pkg/common/errorreport"
	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	"github.com/GoogleCloudPlatform/khi/pkg/parameters"
	privateparameters "github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	"github.com/google/go-cmp/cmp"
)

func TestPrivateIAMTokenInitializer(t *testing.T) {
	testCases := []struct {
		name           string
		inspectionMode bool
		iamToken       string
		fixedProject   string
		gaLabels       string
		wantInjector   bool
	}{
		{
			name:           "inspection mode enabled with fixed project and IAM token",
			inspectionMode: true,
			iamToken:       "test-token",
			fixedProject:   "test-project",
			gaLabels:       "key1=val1,key2=val2",
			wantInjector:   true,
		},
		{
			name:           "inspection mode disabled",
			inspectionMode: false,
			wantInjector:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errorreport.DefaultErrorReporter = errorreport.NewReporter(&errorreport.ConsoleReportWriter{})
			privateparameters.Private.InspectionMode = &tc.inspectionMode
			privateparameters.Private.IAMToken = &tc.iamToken
			privateparameters.Private.GALabels = &tc.gaLabels
			parameters.Auth.FixedProjectID = &tc.fixedProject

			engine := coreinit.NewEngine(context.Background())
			ctx := engine.Context()

			if err := PrivateIAMTokenInitializer.Init(ctx); err != nil {
				t.Fatalf("PrivateIAMTokenInitializer.Init() failed: %v", err)
			}

			injector, found := coreinit.Get(ctx, IAMTokenInjectorKey)
			if found != tc.wantInjector {
				t.Errorf("IAMTokenInjectorKey presence mismatch (-want %v, +got %v)", tc.wantInjector, found)
			}
			if tc.wantInjector && injector == nil {
				t.Errorf("expected injector to be non-nil")
			}

			if diff := cmp.Diff(PrivateIAMTokenInitializer.ID, InitializerIDPrivateIAMToken); diff != "" {
				t.Errorf("Initializer ID mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
