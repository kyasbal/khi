package privatecommon_impl

import (
	"context"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
)

// APICallOptionsInjectorTask is the default implementation to provide the CallOptionInjector.
// Each APIClient use must call this injector method before to supply parameters correctly.
var APICallOptionsInjectorTask = inspectiontaskbase.NewInspectionTask(
	privatecommon_contract.APIClientCallOptionsInjectorTaskOverrideID,
	[]taskid.UntypedTaskReference{},
	func(ctx context.Context, taskMode inspectioncore_contract.InspectionTaskModeType) (*googlecloud.CallOptionInjector, error) {
		options, _ := khictx.GetValue(ctx, privatecommon_contract.APICallOptionsInjectorContextKey)
		if options == nil { // It's OK if the context value wasn't provided then it's not for inspection mode.
			return googlecloud.NewCallOptionInjector(), nil
		}
		return googlecloud.NewCallOptionInjector(*options...), nil
	},
	coretask.WithSelectionPriority(googlecloudcommon_contract.DefaultAPIClientOptionTasksPriority+1),
)
