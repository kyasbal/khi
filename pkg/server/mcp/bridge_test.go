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
	"context"
	"testing"
	"time"

	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestBridgeConns(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "forwards messages between two in-memory transports",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			taskServer, err := coreinspection.NewServer(nil)
			if err != nil {
				t.Fatalf("NewServer failed: %v", err)
			}
			mcpServer := NewServer(taskServer)

			// tServer <-> tBridge1
			tServer, tBridge1 := mcp.NewInMemoryTransports()
			// tBridge2 <-> tClient
			tBridge2, tClient := mcp.NewInMemoryTransports()

			go func() {
				_ = mcpServer.MCPServer().Run(ctx, tServer)
			}()

			conn1, err := tBridge1.Connect(ctx)
			if err != nil {
				t.Fatalf("tBridge1.Connect failed: %v", err)
			}
			defer conn1.Close()

			conn2, err := tBridge2.Connect(ctx)
			if err != nil {
				t.Fatalf("tBridge2.Connect failed: %v", err)
			}
			defer conn2.Close()

			go func() {
				_ = BridgeConns(ctx, conn1, conn2)
			}()

			client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
			session, err := client.Connect(ctx, tClient, nil)
			if err != nil {
				t.Fatalf("client.Connect failed: %v", err)
			}
			defer session.Close()

			res, err := session.CallTool(ctx, &mcp.CallToolParams{
				Name: "list_inspections",
			})
			if err != nil {
				t.Fatalf("CallTool failed: %v", err)
			}
			if res.IsError {
				t.Fatalf("CallTool returned error result")
			}
		})
	}
}
