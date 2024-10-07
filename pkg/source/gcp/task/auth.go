package task

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/env"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/api/iamtoken"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api/accesstoken"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const GCPApiClientTaskId = GCPPrefix + "api-client"

var GCPApiClientTask = task.NewProcessorTask(GCPApiClientTaskId,
	[]string{
		EnvQuotaProjectIdVariableName,
	},
	func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
		quotaProject, err := env.GetEnvironmentVariableFromTaskVariables(ctx, EnvQuotaProjectIdVariableName, v)
		if err != nil {
			return nil, err
		}
		return api.NewGCPClient(accesstoken.DefaultAccessTokenStore, iamtoken.DefaultIAMTokenStore, quotaProject.Value)
	})

func GetGCPApiClientFromTaskVariable(v *task.VariableSet) (api.GCPClient, error) {
	return task.GetTypedVariableFromTaskVariable[api.GCPClient](v, GCPApiClientTaskId, nil)
}
