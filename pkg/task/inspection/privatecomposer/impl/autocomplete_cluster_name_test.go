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

package privatecomposer_impl

import (
	"context"
	"fmt"
	"testing"
	"time"

	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	googlecloudclustercomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudclustercomposer/contract"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecomposer_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecomposer/contract"
	"github.com/google/go-cmp/cmp"
)

type mockComposerClusterFinder struct {
	clusterMapping map[string][]string // {projectID}/{location}/{environment} -> []string
	wantError      bool
}

// GetGKEClusterNames implements googlecloudclustercomposer_contract.ComposerEnvironmentClusterFinder.
func (m *mockComposerClusterFinder) GetGKEClusterNames(ctx context.Context, projectID, location, environment string, startTime, endTime time.Time) ([]string, error) {
	if m.wantError {
		return nil, fmt.Errorf("test error")
	}
	key := fmt.Sprintf("%s/%s/%s", projectID, location, environment)
	if clusterNames, ok := m.clusterMapping[key]; ok {
		if len(clusterNames) == 0 {
			return nil, googlecloudclustercomposer_contract.ErrEnvironmentClusterNotFound
		}
		return clusterNames, nil
	}
	return nil, googlecloudclustercomposer_contract.ErrEnvironmentClusterNotFound
}

var _ googlecloudclustercomposer_contract.ComposerEnvironmentClusterFinder = (*mockComposerClusterFinder)(nil)

func TestAutocompleteComposerClusterIdentityTask(t *testing.T) {
	testCases := []struct {
		desc           string
		clusterMapping map[string][]string
		finderError    bool
		tenantIDs      []string
		environments   []string
		locations      []string
		want           []*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]
	}{
		{
			desc:           "tenant project id is empty",
			clusterMapping: map[string][]string{},
			finderError:    false,
			tenantIDs:      []string{""},
			environments:   []string{"env1"},
			locations:      []string{"us-central1"},
			want: []*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]{{
				Values: []googlecloudk8scommon_contract.GoogleCloudClusterIdentity{},
				Error:  "Project ID or Composer environment name is empty",
			}},
		},
		{
			desc:           "environment name is empty",
			clusterMapping: map[string][]string{},
			finderError:    false,
			tenantIDs:      []string{"foo-tenant-tp"},
			environments:   []string{""},
			locations:      []string{"us-central1"},
			want: []*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]{{
				Values: []googlecloudk8scommon_contract.GoogleCloudClusterIdentity{},
				Error:  "Project ID or Composer environment name is empty",
			}},
		},
		{
			desc:           "location is empty returns hint",
			clusterMapping: map[string][]string{},
			finderError:    false,
			tenantIDs:      []string{"foo-tenant-tp"},
			environments:   []string{"env1"},
			locations:      []string{""},
			want: []*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]{{
				Values: []googlecloudk8scommon_contract.GoogleCloudClusterIdentity{},
				Error:  "",
				Hint:   "Cluster names are suggested after the location is provided.",
			}},
		},
		{
			desc:           "using cache with multiple clusters",
			clusterMapping: map[string][]string{"foo-tenant-tp/us-central1/env1": {"cluster1", "cluster2"}},
			finderError:    false,
			tenantIDs:      []string{"foo-tenant-tp", "foo-tenant-tp"},
			environments:   []string{"env1", "env1"},
			locations:      []string{"us-central1", "us-central1"},
			want: []*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]{
				{
					Values: []googlecloudk8scommon_contract.GoogleCloudClusterIdentity{
						{
							ClusterName:  "cluster1",
							ProjectID:    "foo-tenant-tp",
							Location:     "us-central1",
							PrefixPolicy: googlecloudk8scommon_contract.ClusterPrefixPolicy{},
						},
						{
							ClusterName:  "cluster2",
							ProjectID:    "foo-tenant-tp",
							Location:     "us-central1",
							PrefixPolicy: googlecloudk8scommon_contract.ClusterPrefixPolicy{},
						},
					},
				},
				{
					Values: []googlecloudk8scommon_contract.GoogleCloudClusterIdentity{
						{
							ClusterName:  "cluster1",
							ProjectID:    "foo-tenant-tp",
							Location:     "us-central1",
							PrefixPolicy: googlecloudk8scommon_contract.ClusterPrefixPolicy{},
						},
						{
							ClusterName:  "cluster2",
							ProjectID:    "foo-tenant-tp",
							Location:     "us-central1",
							PrefixPolicy: googlecloudk8scommon_contract.ClusterPrefixPolicy{},
						},
					},
				},
			},
		},
		{
			desc:           "with error",
			clusterMapping: map[string][]string{},
			finderError:    true,
			tenantIDs:      []string{"foo-tenant-tp"},
			environments:   []string{"env1"},
			locations:      []string{"us-central1"},
			want: []*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]{{
				Values: []googlecloudk8scommon_contract.GoogleCloudClusterIdentity{},
				Error:  "Failed to fetch the list GKE cluster. Please confirm if the Project ID is correct, or retry later",
			}},
		},
		{
			desc:           "environment not found",
			clusterMapping: map[string][]string{},
			finderError:    false,
			tenantIDs:      []string{"foo-tenant-tp"},
			environments:   []string{"non-existent-env"},
			locations:      []string{"us-central1"},
			want: []*inspectioncore_contract.AutocompleteResult[googlecloudk8scommon_contract.GoogleCloudClusterIdentity]{{
				Values: []googlecloudk8scommon_contract.GoogleCloudClusterIdentity{},
				Error: `Not found. It works for the clusters existed in the past but make sure the cluster name is right if you believe the cluster should be there.
Note: If you want to inspect Managed Airflow 2, you should select another inspection type.`,
			}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(t.Context())

			mockComposerClusterFinderInput := tasktest.NewTaskDependencyValuePair[googlecloudclustercomposer_contract.ComposerEnvironmentClusterFinder](
				googlecloudclustercomposer_contract.ComposerEnvironmentClusterFinderTaskID.Ref(),
				&mockComposerClusterFinder{
					clusterMapping: tc.clusterMapping,
					wantError:      tc.finderError,
				},
			)

			for i := 0; i < len(tc.tenantIDs); i++ {
				tenantIDInput := tasktest.NewTaskDependencyValuePair(privatecomposer_contract.InputComposerTenantProjectIdTaskID.Ref(), tc.tenantIDs[i])
				environmentNameInput := tasktest.NewTaskDependencyValuePair(googlecloudclustercomposer_contract.InputComposerEnvironmentNameTaskID.Ref(), tc.environments[i])
				locationInput := tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputLocationsTaskID.Ref(), tc.locations[i])
				startTimeInput := tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputStartTimeTaskID.Ref(), time.Unix(1700000000, 0))
				endTimeInput := tasktest.NewTaskDependencyValuePair(googlecloudcommon_contract.InputEndTimeTaskID.Ref(), time.Unix(1700003600, 0))
				result, _, err := inspectiontest.RunInspectionTask(ctx, AutocompleteComposerClusterIdentityTask, inspectioncore_contract.TaskModeDryRun, map[string]any{}, tenantIDInput, environmentNameInput, locationInput, startTimeInput, endTimeInput, mockComposerClusterFinderInput)
				if err != nil {
					t.Fatalf("failed to run inspection task in loop %d: %v", i, err)
				}

				if diff := cmp.Diff(tc.want[i], result); diff != "" {
					t.Errorf("result of AutocompleteComposerClusterIdentityTask mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
