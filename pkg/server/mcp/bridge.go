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
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BridgeConns bridges two bidirectional MCP Connections by forwarding messages in both directions until either closes.
func BridgeConns(ctx context.Context, c1, c2 mcp.Connection) error {
	errCh := make(chan error, 2)

	go func() {
		for {
			msg, err := c1.Read(ctx)
			if err != nil {
				errCh <- err
				return
			}
			if err := c2.Write(ctx, msg); err != nil {
				errCh <- err
				return
			}
		}
	}()

	go func() {
		for {
			msg, err := c2.Read(ctx)
			if err != nil {
				errCh <- err
				return
			}
			if err := c1.Write(ctx, msg); err != nil {
				errCh <- err
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
}

// RunStdioBridge connects to an active KHI server's SSE endpoint and bridges stdin/stdout to it.
func RunStdioBridge(ctx context.Context, serverURL string) error {
	cleanURL := strings.TrimSuffix(serverURL, "/")
	sseEndpoint := cleanURL + "/mcp/sse"

	sseTransport := &mcp.SSEClientTransport{
		Endpoint: sseEndpoint,
	}

	sseConn, err := sseTransport.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to KHI server at %s: %w (please start KHI first)", sseEndpoint, err)
	}
	defer sseConn.Close()

	stdioTransport := &mcp.StdioTransport{}
	stdioConn, err := stdioTransport.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize stdio transport: %w", err)
	}
	defer stdioConn.Close()

	return BridgeConns(ctx, stdioConn, sseConn)
}
