package task

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/api/iamtoken"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api/accesstoken"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const GCPApiClientTaskId = GCPPrefix + "api-client"

var GCPApiClientTask = task.NewProcessorTask(GCPApiClientTaskId,
	[]string{},
	func(ctx context.Context, taskMode int, v *task.VariableSet) (any, error) {
		return api.NewGCPClient(accesstoken.DefaultAccessTokenStore, iamtoken.DefaultIAMTokenStore, *parameters.Auth.QuotaProjectID)
	})

func GetGCPApiClientFromTaskVariable(v *task.VariableSet) (api.GCPClient, error) {
	return task.GetTypedVariableFromTaskVariable[api.GCPClient](v, GCPApiClientTaskId, nil)
}
