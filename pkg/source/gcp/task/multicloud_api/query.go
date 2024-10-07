package multicloud_api

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var MultiCloudAPIQueryTaskId = query.GKEQueryPrefix + "multicloud-api"

func GenerateMultiCloudAPIQuery(clusterNameWithPrefix string) string {
	return fmt.Sprintf(`resource.type="audited_resource"
resource.labels.service="gkemulticloud.googleapis.com"
resource.labels.method:("Update" OR "Create" OR "Delete")
protoPayload.resourceName:"%s"
`, clusterNameWithPrefix)
}

var MultiCloudAPIQueryTask = query.NewQueryGeneratorTask(MultiCloudAPIQueryTaskId, "Multicloud API Logs", enum.LogTypeMulticloudAPI, []string{
	gcp_task.InputClusterName,
}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
	clusterName, err := gcp_task.GetInputClusterNameFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	return []string{GenerateMultiCloudAPIQuery(clusterName)}, nil
})
