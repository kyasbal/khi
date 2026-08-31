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

package defaultinit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	"github.com/GoogleCloudPlatform/khi/pkg/parameters"
	"github.com/GoogleCloudPlatform/khi/pkg/server/workbench"
	"github.com/gin-gonic/gin"
)

func TestMCPServerInitializer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskServer, err := coreinspection.NewServer(nil)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	testCases := []struct {
		name        string
		jobMode     bool
		basePath    string
		method      string
		requestPath string
		wantHit     bool
	}{
		{
			name:        "registers SSE route at root when not in job mode and empty base path",
			jobMode:     false,
			basePath:    "",
			method:      http.MethodGet,
			requestPath: "/mcp/sse",
			wantHit:     true,
		},
		{
			name:        "registers SSE route with custom base path",
			jobMode:     false,
			basePath:    "/custom/prefix",
			method:      http.MethodGet,
			requestPath: "/custom/prefix/mcp/sse",
			wantHit:     true,
		},
		{
			name:        "skips registration when in job mode",
			jobMode:     true,
			basePath:    "",
			method:      http.MethodGet,
			requestPath: "/mcp/sse",
			wantHit:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			engine := coreinit.NewEngine(context.Background())
			ctx := engine.Context()
			jobParams := &parameters.JobParameters{
				JobMode: &tc.jobMode,
			}
			coreinit.Set(ctx, JobParametersKey, jobParams)
			coreinit.Set(ctx, InspectionTaskServerKey, taskServer)

			indexMgr := workbench.NewInspectionIndexManager(taskServer, t.TempDir())
			wbMgr := workbench.NewWorkbenchManager(taskServer, indexMgr, 5*time.Minute, 0)
			defer wbMgr.Stop()
			coreinit.Set(ctx, WorkbenchManagerKey, wbMgr)

			ginEngine := gin.New()
			var router gin.IRouter = ginEngine
			if tc.basePath != "" {
				router = ginEngine.Group(tc.basePath)
			}

			coreinit.Set(ctx, GinRouterKey, router)
			coreinit.Set(ctx, BasePathKey, tc.basePath)

			if err := MCPServerInitializer.Init(ctx); err != nil {
				t.Fatalf("MCPServerInitializer.Init() failed: %v", err)
			}

			reqCtx, reqCancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer reqCancel()

			req := httptest.NewRequestWithContext(reqCtx, tc.method, tc.requestPath, nil)
			w := httptest.NewRecorder()
			ginEngine.ServeHTTP(w, req)

			if tc.wantHit && w.Code == http.StatusNotFound {
				t.Errorf("expected route %q to be registered, got 404", tc.requestPath)
			}
			if !tc.wantHit && w.Code != http.StatusNotFound {
				t.Errorf("expected route %q not to be registered (404), got %d", tc.requestPath, w.Code)
			}
		})
	}
}
