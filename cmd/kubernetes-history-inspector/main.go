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

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/GoogleCloudPlatform/khi/pkg/common/errorreport"
	coreinit "github.com/GoogleCloudPlatform/khi/pkg/core/init"
	"github.com/GoogleCloudPlatform/khi/pkg/lifecycle"
	"github.com/GoogleCloudPlatform/khi/pkg/server/mcp"

	_ "github.com/GoogleCloudPlatform/khi/pkg/core/init/default"
)

func handleTerminateSignal(engine *coreinit.Engine, exitCh chan<- int) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	s := <-sig
	lifecycle.Default.NotifyTerminate(s)
	if err := engine.Terminate(); err != nil {
		slog.Error("Calling termination hooks on Engine failed", "error", err)
	}

	exitCh <- 128 + int(s.(syscall.Signal))
}

func main() {
	// main() shouldn't have any lines other than this, to prevent calling defer errorreport.CheckAndReportPanic on os.Exit
	os.Exit(run())
}

func run() int {
	defer errorreport.CheckAndReportPanic()

	// Check if running in MCP stdio bridge mode.
	for i, arg := range os.Args[1:] {
		if arg == "mcp" || arg == "--mcp" || arg == "-mcp" {
			serverURL := "http://127.0.0.1:8080"
			if envURL := os.Getenv("KHI_SERVER_URL"); envURL != "" {
				serverURL = envURL
			}
			for j, a := range os.Args[1:] {
				if strings.HasPrefix(a, "--mcp-server-url=") {
					serverURL = strings.TrimPrefix(a, "--mcp-server-url=")
				} else if a == "--mcp-server-url" && j+2 < len(os.Args) {
					serverURL = os.Args[j+2]
				}
			}
			_ = i
			if err := mcp.RunStdioBridge(context.Background(), serverURL); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return 1
			}
			return 0
		}
	}

	engine := coreinit.NewEngine(context.Background())
	defer func() {
		if err := engine.Terminate(); err != nil {
			slog.Error("Calling termination hooks on Engine failed", "error", err)
		}
	}()

	if err := engine.Init(); err != nil {
		slog.Error("Initializing KHI failed", "error", err)
		return 1
	}

	exitCh := make(chan int, 1)
	go handleTerminateSignal(engine, exitCh)

	if err := engine.Run(); err != nil {
		slog.Error("Running KHI failed", "error", err)
		return 1
	}

	select {
	case code := <-exitCh:
		return code
	default:
		if engine.Context().Err() != nil {
			return <-exitCh
		}
		return 0
	}
}
