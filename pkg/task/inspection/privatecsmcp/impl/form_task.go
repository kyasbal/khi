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

package privatecsmcp_impl

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/formtask"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecsmcp_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecsmcp/contract"
)

const formPriority = googlecloudcommon_contract.PriorityForResourceIdentifierGroup + 1000

var projectIdValidator = regexp.MustCompile(`^\s*[0-9a-z\.:\-]+\s*$`)

// InputCSMTenantProjectIDTask defines a text input form task for the CSM Tenant Project ID.
var InputCSMTenantProjectIDTask = formtask.NewTextFormTaskBuilder(
	privatecsmcp_contract.InputCSMTenantProjectIDTaskID,
	formPriority+1,
	"CSM Tenant Project ID",
).
	WithDescription("The project ID of the CSM Tenant where Cloud Run metrics reside.").
	WithDefaultValueFunc(func(ctx context.Context, previousValues []string) (string, error) {
		if len(previousValues) > 0 {
			return previousValues[0], nil
		}
		return "", nil
	}).
	WithValidatingTiming(inspectionmetadata.Blur).
	WithValidator(func(ctx context.Context, value string) (string, error) {
		trimmed := strings.TrimSpace(value)
		if !projectIdValidator.MatchString(trimmed) {
			return "Project ID must match `^[0-9a-z\\.:\\-]+$`", nil
		}
		if !strings.HasSuffix(trimmed, "-tp") {
			return "CSM Tenant Project ID must end with `-tp`", nil
		}
		return "", nil
	}).
	WithConverter(func(ctx context.Context, value string) (string, error) {
		return strings.TrimSpace(value), nil
	}).
	Build()

// AutocompleteCSMCPCloudRunServiceNameTask is a task that provides a list of available Cloud Run service names for autocomplete.
var AutocompleteCSMCPCloudRunServiceNameTask = inspectiontaskbase.NewGlobalCachedTask(
	privatecsmcp_contract.AutocompleteCSMCPCloudRunServiceNameTaskID,
	[]taskid.UntypedTaskReference{
		privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref(),
		googlecloudcommon_contract.InputStartTimeTaskID.Ref(),
		googlecloudcommon_contract.InputEndTimeTaskID.Ref(),
		googlecloudcommon_contract.APIClientFactoryTaskID.Ref(),
		googlecloudcommon_contract.APIClientCallOptionsInjectorTaskID.Ref(),
	},
	func(ctx context.Context, prevValue inspectiontaskbase.CacheableTaskResult[*inspectioncore_contract.AutocompleteResult[string]]) (inspectiontaskbase.CacheableTaskResult[*inspectioncore_contract.AutocompleteResult[string]], error) {
		tenantProjectID := coretask.GetTaskResult(ctx, privatecsmcp_contract.InputCSMTenantProjectIDTaskID.Ref())
		startTime := coretask.GetTaskResult(ctx, googlecloudcommon_contract.InputStartTimeTaskID.Ref())
		endTime := coretask.GetTaskResult(ctx, googlecloudcommon_contract.InputEndTimeTaskID.Ref())
		cf := coretask.GetTaskResult(ctx, googlecloudcommon_contract.APIClientFactoryTaskID.Ref())
		optionInjector := coretask.GetTaskResult(ctx, googlecloudcommon_contract.APIClientCallOptionsInjectorTaskID.Ref())

		currentDigest := fmt.Sprintf("%s-%d-%d", tenantProjectID, startTime.Unix(), endTime.Unix())
		if currentDigest == prevValue.DependencyDigest {
			return prevValue, nil
		}
		if tenantProjectID == "" {
			return inspectiontaskbase.CacheableTaskResult[*inspectioncore_contract.AutocompleteResult[string]]{
				Value: &inspectioncore_contract.AutocompleteResult[string]{
					Values: []string{},
					Error:  "",
					Hint:   "Service names are suggested after the CSM tenant project ID is provided.",
				},
				DependencyDigest: currentDigest,
			}, nil
		}

		errorString := ""
		hintString := ""
		if endTime.Before(time.Now().Add(-time.Hour * 24 * 30 * 24)) {
			hintString = "The end time is more than 24 months ago. Suggested service names may not be complete."
		}

		client, err := cf.MonitoringMetricClient(ctx, googlecloud.Project(tenantProjectID))
		if err != nil {
			return prevValue, fmt.Errorf("failed to create monitoring metric client: %w", err)
		}

		ctx = optionInjector.InjectToCallContext(ctx, googlecloud.Project(tenantProjectID))
		filter := `metric.type="run.googleapis.com/container/instance_count" AND resource.type="cloud_run_revision"`
		services, err := googlecloud.QueryDistinctStringLabelValuesFromMetrics(ctx, client, tenantProjectID, filter, startTime, endTime, "resource.labels.service_name", "service_name")
		if err != nil {
			errorString = err.Error()
		}
		if hintString == "" && errorString == "" && len(services) == 0 {
			hintString = fmt.Sprintf("No service names found between %s and %s. It is highly likely that the time range is incorrect. Please verify the time range, or proceed by manually entering the service name.", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
		}

		sort.Strings(services)
		return inspectiontaskbase.CacheableTaskResult[*inspectioncore_contract.AutocompleteResult[string]]{
			DependencyDigest: currentDigest,
			Value: &inspectioncore_contract.AutocompleteResult[string]{
				Values: services,
				Error:  errorString,
				Hint:   hintString,
			},
		}, nil
	},
	coretask.WithSelectionPriority(500),
)

// InputCSMCPCloudRunServiceNameTask defines a text input form task for the CSM CP Cloud Run Service Name.
var InputCSMCPCloudRunServiceNameTask = formtask.NewTextFormTaskBuilder(
	privatecsmcp_contract.InputCSMCPCloudRunServiceNameTaskID,
	formPriority,
	"CSM CP Cloud Run Service Name",
).
	WithDependencies([]taskid.UntypedTaskReference{
		privatecsmcp_contract.AutocompleteCSMCPCloudRunServiceNameTaskID.Ref(),
	}).
	WithDescription("The name of the Cloud Run service for CSM Control Plane.").
	WithDefaultValueFunc(func(ctx context.Context, previousValues []string) (string, error) {
		result := coretask.GetTaskResult(ctx, privatecsmcp_contract.AutocompleteCSMCPCloudRunServiceNameTaskID.Ref())
		if len(previousValues) > 0 {
			// Pre-fill if exists in autocomplete list, or keep the previous value.
			for _, val := range result.Values {
				if val == previousValues[0] {
					return previousValues[0], nil
				}
			}
			return previousValues[0], nil
		}
		if len(result.Values) > 0 {
			return result.Values[0], nil
		}
		return "", nil
	}).
	WithSuggestionsFunc(func(ctx context.Context, value string, previousValues []string) ([]string, error) {
		result := coretask.GetTaskResult(ctx, privatecsmcp_contract.AutocompleteCSMCPCloudRunServiceNameTaskID.Ref())
		var matched []string
		for _, val := range result.Values {
			if strings.Contains(strings.ToLower(val), strings.ToLower(value)) {
				matched = append(matched, val)
			}
		}
		return matched, nil
	}).
	WithHintFunc(func(ctx context.Context, value string, convertedValue any) (string, inspectionmetadata.ParameterHintType, error) {
		result := coretask.GetTaskResult(ctx, privatecsmcp_contract.AutocompleteCSMCPCloudRunServiceNameTaskID.Ref())
		if result.Error != "" {
			return fmt.Sprintf("Failed to obtain the service list due to error '%s'.", result.Error), inspectionmetadata.Warning, nil
		}
		if result.Hint != "" {
			return result.Hint, inspectionmetadata.Info, nil
		}
		return "", inspectionmetadata.Info, nil
	}).
	WithValidatingTiming(inspectionmetadata.Blur).
	WithValidator(func(ctx context.Context, value string) (string, error) {
		if strings.TrimSpace(value) == "" {
			return "Cloud Run Service Name must not be empty", nil
		}
		return "", nil
	}).
	WithConverter(func(ctx context.Context, value string) (string, error) {
		return strings.TrimSpace(value), nil
	}).
	Build()
