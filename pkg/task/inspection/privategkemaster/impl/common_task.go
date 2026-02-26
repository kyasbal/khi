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

package privategkemaster_impl

import (
	"context"

	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/enum"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudinspectiontypegroup_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudinspectiontypegroup/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privategkemaster_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privategkemaster/contract"
)

// CommonFieldSetReaderTask reads the component name at first to filter logs for specific components in the later tasks.
var CommonFieldSetReaderTask = inspectiontaskbase.NewFieldSetReadTask(privategkemaster_contract.CommonFieldSetReaderTaskID,
	privategkemaster_contract.ListLogEntriesTaskID.Ref(),
	[]log.FieldSetReader{
		privategkemaster_contract.NewGKEMasterLogFieldSetReader(),
		&privategkemaster_contract.GKEMasterCommonFieldSetReader{},
	},
)

var logIngesterTask = inspectiontaskbase.NewLogIngesterTask(
	privategkemaster_contract.LogIngesterTaskID,
	privategkemaster_contract.ListLogEntriesTaskID.Ref(),
)

// TailTask is the task to ensure all logs are processed.
var TailTask = inspectiontaskbase.NewInspectionTask(privategkemaster_contract.TailTaskID,
	[]taskid.UntypedTaskReference{
		privategkemaster_contract.SchedulerLogToTimelineMapperTaskID.Ref(),
		privategkemaster_contract.ControllerManagerLogToTimelineMapperTaskID.Ref(),
		privategkemaster_contract.OtherLogToTimelineMapperTaskID.Ref(),
		privategkemaster_contract.KubeletLogLogToTimelineMapperTaskID.Ref(),
		privategkemaster_contract.ContainerdLogLogToTimelineMapperTaskID.Ref(),
	},
	func(ctx context.Context, taskMode inspectioncore_contract.InspectionTaskModeType) (struct{}, error) {
		return struct{}{}, nil
	},
	inspectioncore_contract.FeatureTaskLabel(
		"GKE Master Logs(PRIVATE)",
		`GKE KCP logs from the tenant project. You may need to request access via AoD to use this feature. Please check go/khi-master-log for more details.`,
		enum.LogTypeControlPlaneComponent,
		20000,
		false,
		googlecloudinspectiontypegroup_contract.GCPK8sClusterInspectionTypes...,
	),
)
