package privatecommon_impl

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudinspectiontypegroup_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudinspectiontypegroup/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
)

var FilenameHeaderMetadataGeneratorTask = inspectiontaskbase.NewInspectionTask(privatecommon_contract.FileNameHeaderMetadataGeneratorTask, []taskid.UntypedTaskReference{
	privatecommon_contract.JustificationFormTaskID.Ref(),
	googlecloudk8scommon_contract.InputClusterNameTaskID.Ref(),
	googlecloudcommon_contract.InputEndTimeTaskID.Ref(),
	googlecloudcommon_contract.InputStartTimeTaskID.Ref(),
}, func(ctx context.Context, taskMode inspectioncore_contract.InspectionTaskModeType) (struct{}, error) {
	metadataSet := khictx.MustGetValue(ctx, inspectioncore_contract.InspectionRunMetadata)
	header := typedmap.GetOrDefault(metadataSet, inspectionmetadata.HeaderMetadataKey, &inspectionmetadata.HeaderMetadata{})

	clusterName := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.InputClusterNameTaskID.Ref())
	endTime := coretask.GetTaskResult(ctx, googlecloudcommon_contract.InputEndTimeTaskID.Ref())
	startTime := coretask.GetTaskResult(ctx, googlecloudcommon_contract.InputStartTimeTaskID.Ref())
	justification := coretask.GetTaskResult(ctx, privatecommon_contract.JustificationFormTaskID.Ref())

	header.SuggestedFileName = fmt.Sprintf("%s-%s-%s-%s.khi", justification, clusterName, endTime.Format("01021504"), startTime.Format("01021504"))
	return struct{}{}, nil
},
	inspectioncore_contract.InspectionTypeLabel(googlecloudinspectiontypegroup_contract.GCPK8sClusterInspectionTypes...),
	inspectioncore_contract.NewRequiredTaskLabel())
