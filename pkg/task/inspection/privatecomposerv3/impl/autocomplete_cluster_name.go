// Copyright 2025 Google LLC
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

package privatecomposerv3_impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposerv3_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposerv3/contract"
)

var AutocompleteComposerClusterNamesTask = inspectiontaskbase.NewCachedTask(taskid.NewImplementationID(googlecloudk8scommon_contract.AutocompleteClusterIdentityTaskID.Ref(), privatecomposerv3_contract.InspectionTypeId), []taskid.UntypedTaskReference{
	googlecloudk8scommon_contract.ClusterNamePrefixTaskRef,
	privatecomposerv3_contract.InputComposerV3TenantProjectIdTaskID.Ref(),
	googlecloudcommon_contract.InputStartTimeTaskID.Ref(),
	googlecloudcommon_contract.InputEndTimeTaskID.Ref(),
	googlecloudk8scommon_contract.AutocompleteMetricsK8sContainerTaskID.Ref(),
	googlecloudcommon_contract.APIClientFactoryTaskID.Ref(),
	googlecloudcommon_contract.APIClientCallOptionsInjectorTaskID.Ref(),
}, func(ctx context.Context, prevValue inspectiontaskbase.CacheableTaskResult[*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]]) (inspectiontaskbase.CacheableTaskResult[*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]], error) {
	clusterNamePrefix := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterNamePrefixTaskRef)
	projectID := coretask.GetTaskResult(ctx, privatecomposerv3_contract.InputComposerV3TenantProjectIdTaskID.Ref())
	startTime := coretask.GetTaskResult(ctx, googlecloudcommon_contract.InputStartTimeTaskID.Ref())
	endTime := coretask.GetTaskResult(ctx, googlecloudcommon_contract.InputEndTimeTaskID.Ref())
	metricsType := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.AutocompleteMetricsK8sContainerTaskID.Ref())
	cf := coretask.GetTaskResult(ctx, googlecloudcommon_contract.APIClientFactoryTaskID.Ref())
	optionInjector := coretask.GetTaskResult(ctx, googlecloudcommon_contract.APIClientCallOptionsInjectorTaskID.Ref())

	currentDigest := fmt.Sprintf("%s-%s-%d-%d", clusterNamePrefix, projectID, startTime.Unix(), endTime.Unix())
	if currentDigest == prevValue.DependencyDigest {
		return prevValue, nil
	}
	if projectID == "" {
		return inspectiontaskbase.CacheableTaskResult[*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]]{
			Value: &inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]{
				Values: []googlecloudk8scommon_contract.GoogleCloudClusterIdentity{},
				Error:  "",
				Hint:   "Cluster names are suggested after the tenant project ID is provided.",
			},
			DependencyDigest: currentDigest,
		}, nil
	}

	errorString := ""
	hintString := ""
	if endTime.Before(time.Now().Add(-time.Hour * 24 * 30 * 24)) {
		hintString = "The end time is more than 24 months ago. Suggested cluster names may not be complete."
	}

	client, err := cf.MonitoringMetricClient(ctx, googlecloud.Project(projectID))
	if err != nil {
		return prevValue, fmt.Errorf("failed to create monitoring metric client: %w", err)
	}
	defer client.Close()

	ctx = optionInjector.InjectToCallContext(ctx, googlecloud.Project(projectID))
	filter := fmt.Sprintf(`metric.type="%s" AND resource.type="k8s_container"`, metricsType)
	metricsLabels, err := googlecloud.QueryResourceLabelsFromMetrics(ctx, client, projectID, filter, startTime, endTime, []string{"resource.label.cluster_name", "resource.label.location"})
	if err != nil {
		errorString = err.Error()
	}

	var filteredClusters []map[string]string
	for _, labels := range metricsLabels {
		clusterName := labels["cluster_name"]
		if clusterNamePrefix == "" {
			if !strings.Contains(clusterName, "/") {
				filteredClusters = append(filteredClusters, labels)
			}
		} else if strings.HasPrefix(clusterName, clusterNamePrefix) {
			labels["cluster_name"] = strings.TrimPrefix(clusterName, clusterNamePrefix)
			filteredClusters = append(filteredClusters, labels)
		}
	}

	if hintString == "" && errorString == "" && len(filteredClusters) == 0 {
		hintString = fmt.Sprintf("No cluster names found between %s and %s. It is highly likely that the time range is incorrect. Please verify the time range, or proceed by manually entering the cluster name.", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	}

	identities := make([]googlecloudk8scommon_contract.GoogleCloudClusterIdentity, len(filteredClusters))
	for i, labels := range filteredClusters {
		identities[i] = googlecloudk8scommon_contract.GoogleCloudClusterIdentity{
			ProjectID:         projectID,
			ClusterTypePrefix: clusterNamePrefix,
			ClusterName:       labels["cluster_name"],
			Location:          labels["location"],
		}
	}

	return inspectiontaskbase.CacheableTaskResult[*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]]{
		DependencyDigest: currentDigest,
		Value: &inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]{
			Values: identities,
			Error:  errorString,
			Hint:   hintString,
		},
	}, nil
}, inspectioncore_contract.InspectionTypeLabel(privatecomposerv3_contract.InspectionTypeId),
	coretask.WithSelectionPriority(2000),
)
