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

	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/gcpqueryutil"
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	khifilev6 "github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudlogk8scontrolplane_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8scontrolplane/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privategkemaster_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privategkemaster/contract"
)

// CommonFieldSetReaderTask reads the component name at first to filter logs for specific components in the later tasks.
var CommonFieldSetReaderTask = inspectiontaskbase.NewFieldSetReadTask(privategkemaster_contract.CommonFieldSetReaderTaskID,
	privategkemaster_contract.ListLogEntriesTaskID.Ref(),
	[]log.FieldSetReader{
		privategkemaster_contract.NewGKEMasterLogFieldSetReader(),
		&privategkemaster_contract.GKEMasterCommonFieldSetReader{},
		&gcpqueryutil.GCPCommonFieldSetReader{},
		&googlecloudcommon_contract.GCPDefaultSeverityFieldSetReader{},
	},
)

// PrivateGKEMasterLogIngester is a log ingester for private GKE master logs.
type PrivateGKEMasterLogIngester struct{}

// RawLogTask implements inspectiontaskbase.LogIngester.
func (i *PrivateGKEMasterLogIngester) RawLogTask() taskid.TaskReference[[]*log.Log] {
	return privategkemaster_contract.CommonFieldSetReaderTaskID.Ref()
}

// Dependencies implements inspectiontaskbase.LogIngester.
func (i *PrivateGKEMasterLogIngester) Dependencies() []taskid.UntypedTaskReference {
	return []taskid.UntypedTaskReference{}
}

// ProcessLog implements inspectiontaskbase.LogIngester.
func (i *PrivateGKEMasterLogIngester) ProcessLog(ctx context.Context, l *log.Log) (*khifilev6.LogChangeSet, error) {
	cs, err := khifilev6.NewLogChangeSet(l)
	if err != nil {
		return nil, err
	}
	cs.SetLogType(googlecloudlogk8scontrolplane_contract.LogTypeControlPlaneComponent)

	if commonFS, err := log.GetFieldSet(l, &log.CommonFieldSet{}); err == nil {
		cs.SetTimestamp(commonFS.Timestamp)
	}

	if severityFS, err := log.GetFieldSet(l, &inspectioncore_contract.DefaultSeverityFieldSet{}); err == nil {
		cs.SetSeverity(severityFS.Severity)
	}

	if msgFS, err := log.GetFieldSet(l, &googlecloudlogk8scontrolplane_contract.K8sControlplaneCommonMessageFieldSet{}); err == nil {
		cs.SetSummary(msgFS.Message)
	}

	return cs, nil
}

var _ inspectiontaskbase.LogIngester = (*PrivateGKEMasterLogIngester)(nil)

var logIngesterTask = inspectiontaskbase.NewLogIngesterTask(
	privategkemaster_contract.LogIngesterTaskID,
	&PrivateGKEMasterLogIngester{},
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
		20000,
		false,
	),
)
