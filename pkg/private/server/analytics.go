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

package server

import (
	"context"

	"connectrpc.com/connect"
	apiv1 "github.com/GoogleCloudPlatform/khi/pkg/generated/api/v1"
	"github.com/GoogleCloudPlatform/khi/pkg/generated/api/v1/apiv1connect"
	"github.com/GoogleCloudPlatform/khi/pkg/private/analytics/types"
)

// EventReporter defines the interface for reporting analytics events.
type EventReporter interface {
	ReportEvent(event types.AnalyticsEvent, metadata map[string]any)
}

// PrivateAnalyticsServer implements the PrivateAnalyticsService Connect-RPC handler.
type PrivateAnalyticsServer struct {
	reporter EventReporter
}

var _ apiv1connect.PrivateAnalyticsServiceHandler = (*PrivateAnalyticsServer)(nil)

// NewPrivateAnalyticsServer returns a new instance of PrivateAnalyticsServer.
func NewPrivateAnalyticsServer(reporter EventReporter) *PrivateAnalyticsServer {
	return &PrivateAnalyticsServer{
		reporter: reporter,
	}
}

// ReportActivity reports user activity to the analytics backend.
func (s *PrivateAnalyticsServer) ReportActivity(
	ctx context.Context,
	req *connect.Request[apiv1.ReportActivityRequest],
) (*connect.Response[apiv1.ReportActivityResponse], error) {
	msg := req.Msg
	var event types.AnalyticsEvent
	metadata := map[string]any{}

	switch payload := msg.GetPayload().(type) {
	case *apiv1.ReportActivityRequest_Init:
		event = types.AnalyticsEventFrontendInit
		metadata["pageType"] = payload.Init.GetPageType()
		if frontendVersion := payload.Init.GetFrontendVersion(); frontendVersion != "" {
			metadata["frontendVersion"] = frontendVersion
		}
	case *apiv1.ReportActivityRequest_Inspect:
		event = types.AnalyticsEventFrontendInspect
	case *apiv1.ReportActivityRequest_OpenInspectionData:
		event = types.AnalyticsEventFrontendOpenInspectionData
		metadata["inspectionDataHash"] = payload.OpenInspectionData.GetInspectionDataHash()
		metadata["openId"] = payload.OpenInspectionData.GetOpenId()
		metadata["logLength"] = payload.OpenInspectionData.GetLogLength()
		metadata["decompressedTextBufferLength"] = payload.OpenInspectionData.GetDecompressedTextBufferLength()
		metadata["revisionCount"] = payload.OpenInspectionData.GetRevisionCount()
		metadata["eventCount"] = payload.OpenInspectionData.GetEventCount()
	}

	if event != "" {
		s.reporter.ReportEvent(event, metadata)
	}

	return connect.NewResponse(&apiv1.ReportActivityResponse{}), nil
}
