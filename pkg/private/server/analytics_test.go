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
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	apiv1 "github.com/GoogleCloudPlatform/khi/pkg/generated/api/v1"
	"github.com/GoogleCloudPlatform/khi/pkg/generated/api/v1/apiv1connect"
	"github.com/GoogleCloudPlatform/khi/pkg/private/analytics/types"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

type mockEventReporter struct {
	reportedEvents   []types.AnalyticsEvent
	reportedMetadata []map[string]any
}

var _ EventReporter = (*mockEventReporter)(nil)

func (m *mockEventReporter) ReportEvent(event types.AnalyticsEvent, metadata map[string]any) {
	m.reportedEvents = append(m.reportedEvents, event)
	m.reportedMetadata = append(m.reportedMetadata, metadata)
}

func TestPrivateAnalyticsServer_ReportActivity(t *testing.T) {
	testCases := []struct {
		name         string
		request      *apiv1.ReportActivityRequest
		wantRes      *apiv1.ReportActivityResponse
		wantEvent    types.AnalyticsEvent
		wantMetadata map[string]any
		wantErr      bool
	}{
		{
			name: "Init activity payload with frontend version",
			request: &apiv1.ReportActivityRequest{
				Payload: &apiv1.ReportActivityRequest_Init{
					Init: &apiv1.InitActivityPayload{
						PageType:        proto.String("MAIN"),
						FrontendVersion: proto.String("1.0.0"),
					},
				},
			},
			wantRes:   &apiv1.ReportActivityResponse{},
			wantEvent: types.AnalyticsEventFrontendInit,
			wantMetadata: map[string]any{
				"pageType":        "MAIN",
				"frontendVersion": "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "Init activity payload without frontend version",
			request: &apiv1.ReportActivityRequest{
				Payload: &apiv1.ReportActivityRequest_Init{
					Init: &apiv1.InitActivityPayload{
						PageType: proto.String("MAIN"),
					},
				},
			},
			wantRes:   &apiv1.ReportActivityResponse{},
			wantEvent: types.AnalyticsEventFrontendInit,
			wantMetadata: map[string]any{
				"pageType": "MAIN",
			},
			wantErr: false,
		},
		{
			name: "Init activity payload with empty string frontend version",
			request: &apiv1.ReportActivityRequest{
				Payload: &apiv1.ReportActivityRequest_Init{
					Init: &apiv1.InitActivityPayload{
						PageType:        proto.String("MAIN"),
						FrontendVersion: proto.String(""),
					},
				},
			},
			wantRes:   &apiv1.ReportActivityResponse{},
			wantEvent: types.AnalyticsEventFrontendInit,
			wantMetadata: map[string]any{
				"pageType": "MAIN",
			},
			wantErr: false,
		},
		{
			name: "Inspect activity payload",
			request: &apiv1.ReportActivityRequest{
				Payload: &apiv1.ReportActivityRequest_Inspect{
					Inspect: &apiv1.InspectActivityPayload{},
				},
			},
			wantRes:      &apiv1.ReportActivityResponse{},
			wantEvent:    types.AnalyticsEventFrontendInspect,
			wantMetadata: map[string]any{},
			wantErr:      false,
		},
		{
			name: "OpenInspectionData activity payload",
			request: &apiv1.ReportActivityRequest{
				Payload: &apiv1.ReportActivityRequest_OpenInspectionData{
					OpenInspectionData: &apiv1.OpenInspectionDataActivityPayload{
						InspectionDataHash:           proto.String("hash-123"),
						OpenId:                       proto.String("open-456"),
						LogLength:                    proto.Int64(1000),
						DecompressedTextBufferLength: proto.Int64(50000),
						RevisionCount:                proto.Int64(10),
						EventCount:                   proto.Int64(25),
					},
				},
			},
			wantRes:   &apiv1.ReportActivityResponse{},
			wantEvent: types.AnalyticsEventFrontendOpenInspectionData,
			wantMetadata: map[string]any{
				"inspectionDataHash":           "hash-123",
				"openId":                       "open-456",
				"logLength":                    int64(1000),
				"decompressedTextBufferLength": int64(50000),
				"revisionCount":                int64(10),
				"eventCount":                   int64(25),
			},
			wantErr: false,
		},
		{
			name:      "Nil payload",
			request:   &apiv1.ReportActivityRequest{},
			wantRes:   &apiv1.ReportActivityResponse{},
			wantEvent: "",
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reporter := &mockEventReporter{}
			server := NewPrivateAnalyticsServer(reporter)
			path, handler := apiv1connect.NewPrivateAnalyticsServiceHandler(server)

			mux := http.NewServeMux()
			mux.Handle(path, handler)
			ts := httptest.NewServer(mux)
			defer ts.Close()

			client := apiv1connect.NewPrivateAnalyticsServiceClient(ts.Client(), ts.URL)
			res, err := client.ReportActivity(context.Background(), connect.NewRequest(tc.request))

			if (err != nil) != tc.wantErr {
				t.Fatalf("ReportActivity() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if !tc.wantErr {
				if diff := cmp.Diff(tc.wantRes, res.Msg, protocmp.Transform()); diff != "" {
					t.Errorf("ReportActivity() response mismatch (-want +got):\n%s", diff)
				}
				if tc.wantEvent == "" {
					if len(reporter.reportedEvents) != 0 {
						t.Errorf("ReportActivity() reported unexpected events: %v", reporter.reportedEvents)
					}
				} else {
					if len(reporter.reportedEvents) != 1 {
						t.Fatalf("ReportActivity() reported events count = %d, want 1", len(reporter.reportedEvents))
					}
					if reporter.reportedEvents[0] != tc.wantEvent {
						t.Errorf("ReportActivity() event = %v, want %v", reporter.reportedEvents[0], tc.wantEvent)
					}
					if diff := cmp.Diff(tc.wantMetadata, reporter.reportedMetadata[0]); diff != "" {
						t.Errorf("ReportActivity() metadata mismatch (-want +got):\n%s", diff)
					}
				}
			}
		})
	}
}
