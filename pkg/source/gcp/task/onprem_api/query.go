package onprem_api

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var OnPremCloudAPIQueryTaskId = query.GKEQueryPrefix + "onprem-api"

func GenerateOnPremAPIQuery(clusterNameWithPrefix string) string {
	return fmt.Sprintf(`resource.type="audited_resource"
resource.labels.service="gkeonprem.googleapis.com"
resource.labels.method:("Update" OR "Create" OR "Delete" OR "Enroll" OR "Unenroll")
protoPayload.resourceName:"%s"
`, clusterNameWithPrefix)
}

var OnPremAPIQueryTask = query.NewQueryGeneratorTask(OnPremCloudAPIQueryTaskId, "OnPrem API Logs", enum.LogTypeOnPremAPI, []string{
	gcp_task.InputClusterName,
}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
	clusterName, err := gcp_task.GetInputClusterNameFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	return []string{GenerateOnPremAPIQuery(clusterName)}, nil
})
