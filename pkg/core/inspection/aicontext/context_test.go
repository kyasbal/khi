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

package aicontext

import (
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	"github.com/google/go-cmp/cmp"
)

func TestContextHelpers(t *testing.T) {
	testCases := []struct {
		name         string
		setupCtx     func(aiMeta *inspectionmetadata.AIContextMetadata) context.Context
		callHelpers  func(ctx context.Context)
		wantSections []*inspectionmetadata.AIContextSection
	}{
		{
			name: "records values into ai metadata from context",
			setupCtx: func(aiMeta *inspectionmetadata.AIContextMetadata) context.Context {
				writableMap := typedmap.NewTypedMap()
				typedmap.Set(writableMap, inspectionmetadata.AIContextMetadataKey, aiMeta)
				return khictx.WithValue(context.Background(), inspectioncore_contract.InspectionRunMetadata, writableMap.AsReadonly())
			},
			callHelpers: func(ctx context.Context) {
				SetPriority(ctx, "Test Section", 50)
				SetProperty(ctx, "Test Section", "Key", "Val")
				AddToSetProperty(ctx, "Test Section", "Items", "Item1")
				AddToSetProperty(ctx, "Test Section", "Items", "Item2")
				AppendSummaryMarkdown(ctx, "Test Section", "Markdown content")
			},
			wantSections: []*inspectionmetadata.AIContextSection{
				{
					Title:    "Test Section",
					Priority: 50,
					Properties: map[string]string{
						"Key": "Val",
					},
					SetProperties: map[string][]string{
						"Items": {"Item1", "Item2"},
					},
					SummaryMarkdown: "Markdown content",
				},
			},
		},
		{
			name: "does not panic when context has no metadata",
			setupCtx: func(aiMeta *inspectionmetadata.AIContextMetadata) context.Context {
				return context.Background()
			},
			callHelpers: func(ctx context.Context) {
				SetPriority(ctx, "Test", 50)
				SetProperty(ctx, "Test", "Key", "Val")
				AddToSetProperty(ctx, "Test", "Items", "Item")
				AppendSummaryMarkdown(ctx, "Test", "Markdown")
			},
			wantSections: []*inspectionmetadata.AIContextSection{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aiMeta := inspectionmetadata.NewAIContextMetadata()
			ctx := tc.setupCtx(aiMeta)
			tc.callHelpers(ctx)

			got := aiMeta.Sections()
			if diff := cmp.Diff(tc.wantSections, got); diff != "" {
				t.Errorf("AIContext sections mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
