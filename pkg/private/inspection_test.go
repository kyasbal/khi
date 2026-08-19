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

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	defaultinit "github.com/GoogleCloudPlatform/khi/pkg/core/init/default"
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api/iamtoken"
	privateparameters "github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	"github.com/google/go-cmp/cmp"
)

func TestPrivateInspectionInitializer(t *testing.T) {
	testCases := []struct {
		name           string
		inspectionMode bool
	}{
		{
			name:           "inspection mode enabled adds run context options",
			inspectionMode: true,
		},
		{
			name:           "inspection mode disabled",
			inspectionMode: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			privateparameters.Private.InspectionMode = &tc.inspectionMode

			engine := coreinit.NewEngine(context.Background())
			ctx := engine.Context()

			ioconfig := &inspectioncore_contract.IOConfig{}
			server, err := coreinspection.NewServer(ioconfig)
			if err != nil {
				t.Fatalf("failed to create inspection server: %v", err)
			}
			coreinit.Set(ctx, defaultinit.InspectionTaskServerKey, server)

			if tc.inspectionMode {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("test-project"), "test-token")
				coreinit.Set(ctx, IAMTokenInjectorKey, injector)
			}

			if err := PrivateInspectionInitializer.Init(ctx); err != nil {
				t.Fatalf("PrivateInspectionInitializer.Init() failed: %v", err)
			}

			if diff := cmp.Diff(PrivateInspectionInitializer.ID, InitializerIDPrivateInspection); diff != "" {
				t.Errorf("Initializer ID mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
