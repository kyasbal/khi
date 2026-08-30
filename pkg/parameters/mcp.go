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

package parameters

import (
	"os"

	"github.com/GoogleCloudPlatform/khi/pkg/common/flag"
)

// MCP is the global MCPParameters instance.
var MCP *MCPParameters = &MCPParameters{}

// MCPParameters contains parameters for running KHI in MCP stdio client bridge mode.
type MCPParameters struct {
	// Mode indicates whether KHI is invoked in MCP stdio client bridge mode.
	Mode *bool
	// ServerURL is the URL of the running KHI web server to bridge to.
	ServerURL *string
}

// Prepare implements ParameterStore.
func (m *MCPParameters) Prepare() error {
	m.Mode = flag.Bool("mcp", false, "Run KHI in MCP (Model Context Protocol) stdio client bridge mode.", "")
	m.ServerURL = flag.String("mcp-server-url", "http://127.0.0.1:8080", "The URL of the running KHI server when using MCP mode.", "KHI_SERVER_URL")
	return nil
}

// PostProcess implements ParameterStore.
func (m *MCPParameters) PostProcess() error {
	if len(os.Args) > 1 && os.Args[1] == "mcp" {
		*m.Mode = true
	}
	return nil
}

var _ ParameterStore = (*MCPParameters)(nil)
