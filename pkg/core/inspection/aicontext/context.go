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

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
)

// AddToSetProperty adds a string value into a set property under the specified section in the AI context.
func AddToSetProperty(ctx context.Context, sectionTitle, key, value string) {
	if aiMeta := fromContext(ctx); aiMeta != nil {
		aiMeta.AddToSetProperty(sectionTitle, key, value)
	}
}

// SetProperty records or updates a single key-value property under the specified section in the AI context.
func SetProperty(ctx context.Context, sectionTitle, key, value string) {
	if aiMeta := fromContext(ctx); aiMeta != nil {
		aiMeta.SetProperty(sectionTitle, key, value)
	}
}

// AppendSummaryMarkdown appends markdown text under the specified section in the AI context.
func AppendSummaryMarkdown(ctx context.Context, sectionTitle, markdown string) {
	if aiMeta := fromContext(ctx); aiMeta != nil {
		aiMeta.AppendSummaryMarkdown(sectionTitle, markdown)
	}
}

// SetPriority sets the display priority of the specified section in the AI context.
func SetPriority(ctx context.Context, sectionTitle string, priority int32) {
	if aiMeta := fromContext(ctx); aiMeta != nil {
		aiMeta.SetPriority(sectionTitle, priority)
	}
}

// fromContext retrieves the AIContextMetadata instance from the context metadata map, or returns nil if absent.
func fromContext(ctx context.Context) *inspectionmetadata.AIContextMetadata {
	metadataSet, err := khictx.GetValue(ctx, inspectioncore_contract.InspectionRunMetadata)
	if err != nil || metadataSet == nil {
		return nil
	}
	aiMeta, found := typedmap.Get(metadataSet, inspectionmetadata.AIContextMetadataKey)
	if !found {
		return nil
	}
	return aiMeta
}
