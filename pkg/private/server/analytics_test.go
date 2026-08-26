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
	"github.com/GoogleCloudPlatform/khi/pkg/private/analytics"
	"github.com/GoogleCloudPlatform/khi/pkg/private/parameters"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestPrivateAnalyticsServer_ReportActivity(t *testing.T) {
	disableAnalytics := true
	parameters.Private.DisableAnalytics = &disableAnalytics

	testCases := []struct {
		name    string
		request *apiv1.ReportActivityRequest
		wantRes *apiv1.ReportActivityResponse
		wantErr bool
	}{
		{
			name: "Init activity payload",
			request: &apiv1.ReportActivityRequest{
				Payload: &apiv1.ReportActivityRequest_Init{
					Init: &apiv1.InitActivityPayload{
						PageType: proto.String("MAIN"),
					},
				},
			},
			wantRes: &apiv1.ReportActivityResponse{},
			wantErr: false,
		},
		{
			name: "Inspect activity payload",
			request: &apiv1.ReportActivityRequest{
				Payload: &apiv1.ReportActivityRequest_Inspect{
					Inspect: &apiv1.InspectActivityPayload{},
				},
			},
			wantRes: &apiv1.ReportActivityResponse{},
			wantErr: false,
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
			wantRes: &apiv1.ReportActivityResponse{},
			wantErr: false,
		},
		{
			name:    "Nil payload",
			request: &apiv1.ReportActivityRequest{},
			wantRes: &apiv1.ReportActivityResponse{},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewPrivateAnalyticsServer(analytics.NewAnalyticsReporter())
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
			}
		})
	}
}
