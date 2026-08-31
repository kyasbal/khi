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
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	"github.com/GoogleCloudPlatform/khi/pkg/server/mcp"
	"github.com/gin-gonic/gin"
)

var (
	// MCPServerKey stores the MCP Server instance in the init context.
	MCPServerKey = typedmap.NewTypedKey[*mcp.Server]("khi.google.com/init/mcp-server")
)

// InitializerIDMCPServer identifies the Initializer that mounts the MCP Server.
const InitializerIDMCPServer coreinit.InitializerID = "khi.default/mcp-server"

// MCPServerInitializer mounts the MCP Server HTTP handler onto the Gin router.
var MCPServerInitializer = &coreinit.Initializer{
	ID: InitializerIDMCPServer,
	Dependencies: []coreinit.InitializerID{
		InitializerIDGinServer,
		InitializerIDInspectionTaskServer,
		InitializerIDWorkbenchService,
	},
	Before: []coreinit.InitializerID{
		InitializerIDServerRunner,
	},
	Init: func(ctx *coreinit.InitContext) error {
		jobParams := coreinit.MustGet(ctx, JobParametersKey)
		if *jobParams.JobMode {
			return nil
		}
		inspectionServer := coreinit.MustGet(ctx, InspectionTaskServerKey)
		workbenchManager := coreinit.MustGet(ctx, WorkbenchManagerKey)
		router := coreinit.MustGet(ctx, GinRouterKey)
		basePath := coreinit.MustGet(ctx, BasePathKey)

		srv := mcp.NewServer(
			mcp.NewInspectionHandler(inspectionServer),
			mcp.NewWorkbenchHandler(workbenchManager),
		)
		coreinit.Set(ctx, MCPServerKey, srv)

		cleanBasePath := strings.TrimSuffix(basePath, "/")
		var httpHandler http.Handler = srv.HTTPHandler()
		if cleanBasePath != "" {
			httpHandler = http.StripPrefix(cleanBasePath, srv.HTTPHandler())
		}

		router.Any("/mcp/*any", gin.WrapH(httpHandler))
		router.Any("/mcp", gin.WrapH(httpHandler))
		return nil
	},
}

func init() {
	coreinit.RegisterInitializer(MCPServerInitializer)
}
