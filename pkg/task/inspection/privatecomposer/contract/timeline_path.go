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

package privatecomposer_contract

import (
	"context"

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
)

// MustCloudSQLInstanceTimeline returns the timeline path for a specific Cloud SQL database instance under a GCP project.
func MustCloudSQLInstanceTimeline(ctx context.Context, projectID, instanceID string) *khifilev6.TimelinePath {
	if projectID == "" {
		projectID = "unknown"
	}
	if instanceID == "" {
		instanceID = "unknown"
	}
	projectPath := googlecloudcommon_contract.MustGCPProjectTimeline(ctx, projectID)
	builder := khictx.MustGetValue(ctx, inspectioncore_contract.Builder)
	return builder.TimelineAccumulator.GetPath(projectPath, khifilev6.PathSegment{
		Name: instanceID,
		Type: TimelineTypeCloudSQLInstance,
	})
}

// MustCloudSQLLogTimeline returns the timeline path for a specific Cloud SQL database log file under a database instance.
func MustCloudSQLLogTimeline(ctx context.Context, projectID, instanceID, logFileName string) *khifilev6.TimelinePath {
	if logFileName == "" {
		logFileName = "unknown"
	}
	instancePath := MustCloudSQLInstanceTimeline(ctx, projectID, instanceID)
	builder := khictx.MustGetValue(ctx, inspectioncore_contract.Builder)
	return builder.TimelineAccumulator.GetPath(instancePath, khifilev6.PathSegment{
		Name: logFileName,
		Type: TimelineTypeCloudSQLLog,
	})
}
