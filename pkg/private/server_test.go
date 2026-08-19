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
	"github.com/GoogleCloudPlatform/khi/pkg/parameters"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api/iamtoken"
	privateparameters "github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

func TestPrivateServerInitializer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		jobMode        bool
		inspectionMode bool
	}{
		{
			name:           "server mode with inspection mode enabled",
			jobMode:        false,
			inspectionMode: true,
		},
		{
			name:           "server mode with inspection mode disabled",
			jobMode:        false,
			inspectionMode: false,
		},
		{
			name:           "job mode skips route configuration",
			jobMode:        true,
			inspectionMode: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			privateparameters.Private.InspectionMode = &tc.inspectionMode

			engine := coreinit.NewEngine(context.Background())
			ctx := engine.Context()

			jobParams := &parameters.JobParameters{
				JobMode: &tc.jobMode,
			}
			coreinit.Set(ctx, defaultinit.JobParametersKey, jobParams)

			ginEngine := gin.New()
			var router gin.IRouter = ginEngine.Group("/test")
			coreinit.Set(ctx, defaultinit.GinRouterKey, router)

			if tc.inspectionMode {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("test-project"), "test-token")
				coreinit.Set(ctx, IAMTokenInjectorKey, injector)
			}

			if err := PrivateServerInitializer.Init(ctx); err != nil {
				t.Fatalf("PrivateServerInitializer.Init() failed: %v", err)
			}

			if diff := cmp.Diff(PrivateServerInitializer.ID, InitializerIDPrivateServer); diff != "" {
				t.Errorf("Initializer ID mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
