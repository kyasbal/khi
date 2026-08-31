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

package mcp

import (
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DomainHandler defines an MCP feature domain handler that registers its tools and resources.
type DomainHandler interface {
	Register(srv *mcp.Server)
}

// Server encapsulates the MCP server and its transport bindings.
type Server struct {
	mcpServer         *mcp.Server
	streamableHandler *mcp.StreamableHTTPHandler
}

// NewServer creates a new KHI MCP server and registers the provided domain handlers.
func NewServer(handlers ...DomainHandler) *Server {
	impl := &mcp.Implementation{
		Name:    "khi-mcp-server",
		Version: "1.0.0",
	}
	mcpSrv := mcp.NewServer(impl, nil)

	s := &Server{
		mcpServer: mcpSrv,
	}

	for _, handler := range handlers {
		s.RegisterHandler(handler)
	}

	s.streamableHandler = mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return s.mcpServer
	}, nil)

	return s
}

// RegisterHandler registers an additional domain handler onto the MCP server.
func (s *Server) RegisterHandler(handler DomainHandler) {
	if handler != nil {
		handler.Register(s.mcpServer)
	}
}

// MCPServer returns the underlying MCP server instance.
func (s *Server) MCPServer() *mcp.Server {
	return s.mcpServer
}

// HTTPHandler returns an http.Handler that handles Streamable HTTP MCP sessions.
func (s *Server) HTTPHandler() http.Handler {
	return s.streamableHandler
}
