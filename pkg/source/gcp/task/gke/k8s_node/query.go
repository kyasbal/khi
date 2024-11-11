// Copyright 2024 Google LLC
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

package k8s_node

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/query"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

func GenerateK8sNodeLogQuery(projectId string, clusterId string) string {
	return fmt.Sprintf(`resource.type="k8s_node"
-logName="projects/%s/logs/events"
resource.labels.cluster_name="%s"
`, projectId, clusterId)
}

const GKENodeLogQueryTaskId = query.GKEQueryPrefix + "k8s-node"

var GKENodeQueryTask = query.NewQueryGeneratorTask(GKENodeLogQueryTaskId, "Kubernetes node log", enum.LogTypeNode, []string{
	gcp_task.InputProjectIdVariableName,
	gcp_task.InputClusterName,
}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
	clusterName, err := gcp_task.GetInputClusterNameFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	projectId, err := gcp_task.GetInputProjectIdFromTaskVariable(vs)
	if err != nil {
		return []string{}, err
	}
	return []string{GenerateK8sNodeLogQuery(projectId, clusterName)}, nil
})
