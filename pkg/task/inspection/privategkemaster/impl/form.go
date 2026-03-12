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

package privategkemaster_impl

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/formtask"
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/gcpqueryutil"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	googlecloudk8scommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudk8scommon/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
	privategkemaster_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privategkemaster/contract"
)

const priorityForPrivateGKEMasterGroup = googlecloudcommon_contract.FormBasePriority + 20000

var storageScopeValidator = regexp.MustCompile(`storageScope=storage,([^;?]+)`)
var tenantProjectIDValidator = regexp.MustCompile(`resource\.labels\.project_id="([^"]+)"`)

var logNameProjectIDValidator = regexp.MustCompile(`logName="projects/([^/]+)/logs/`)

// InputGKEMasterLogSourceTask is the form task to input the master log source from the link for master cloud logging.
var InputGKEMasterLogSourceTask = formtask.NewTextFormTaskBuilder(privategkemaster_contract.InputGKEMasterLogSourceTaskID, priorityForPrivateGKEMasterGroup+2000, "Master Logs Link").
	WithDescription("The panthenon link to the master logs. Please check go/khi-master-log to know how to get this link.").
	WithDependencies([]taskid.UntypedTaskReference{
		googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref(),
		privatecommon_contract.JustificationFormTaskID.Ref(),
	}).
	WithValidator(validateMasterLogLink).
	WithDefaultValueFunc(defaultMasterLogLink).
	WithConverter(convertMasterLogLink).
	Build()

func defaultMasterLogLink(ctx context.Context, previousValues []string) (string, error) {
	if len(previousValues) > 0 {
		return previousValues[0], nil
	}
	return "", nil
}

func validateMasterLogLink(ctx context.Context, value string) (string, error) {
	justification := coretask.GetTaskResult(ctx, privatecommon_contract.JustificationFormTaskID.Ref())
	identity := coretask.GetTaskResult(ctx, googlecloudk8scommon_contract.ClusterIdentityTaskID.Ref())
	link, err := api.ToGKEAdminLink(googlecloud.Project(identity.ProjectID), justification)
	if err != nil {
		return fmt.Sprintf("Failed to generate GKE Admin link. Please check go/khi-master-log to know how to get this link. \n %s", err.Error()), nil
	}
	if value == "" {
		return fmt.Sprintf("Master Logs Link is required. Open %s and refer to go/khi-master-log", link), nil
	}
	unescapedValue, err := url.PathUnescape(value)
	if err != nil {
		return fmt.Sprintf("Invalid URL format. Open %s and refer to go/khi-master-log", link), nil
	}
	if !storageScopeValidator.MatchString(unescapedValue) {
		return fmt.Sprintf("Master Logs Link must contain `storageScope` parameter. Open %s and refer to go/khi-master-log", link), nil
	}
	if !strings.Contains(unescapedValue, "project=") {
		return fmt.Sprintf("Master Logs Link must contain `project` parameter. Open %s and refer to go/khi-master-log", link), nil
	}
	if !tenantProjectIDValidator.MatchString(unescapedValue) && !logNameProjectIDValidator.MatchString(unescapedValue) {
		return fmt.Sprintf("Master Logs Link must contain `resource.labels.project_id` or `logName` with project ID in query. Open %s and refer to go/khi-master-log", link), nil
	}
	return "", nil
}

func convertMasterLogLink(ctx context.Context, value string) (*privategkemaster_contract.LogSource, error) {
	value, err := url.PathUnescape(value)
	if err != nil {
		return nil, err
	}
	tenantProjectID := ""
	projectIDMatches := tenantProjectIDValidator.FindStringSubmatch(value)
	if len(projectIDMatches) >= 2 {
		tenantProjectID = projectIDMatches[1]
	} else {
		logNameMatches := logNameProjectIDValidator.FindStringSubmatch(value)
		if len(logNameMatches) >= 2 {
			tenantProjectID = logNameMatches[1]
		}
	}

	if tenantProjectID == "" {
		// The validator should catch this, we can ignore this case.
		return nil, nil
	}

	matches := storageScopeValidator.FindStringSubmatch(value)
	if len(matches) < 2 {
		// The validator should catch this, we can ignore this case.
		return nil, nil
	}
	return &privategkemaster_contract.LogSource{
		TenantProjectID:     tenantProjectID,
		LogViewResourceName: matches[1],
	}, nil
}

// InputPrivateGKEMasterComponentNameFilterTask is the form task to filter component names.
var InputPrivateGKEMasterComponentNameFilterTask = formtask.NewSetFormTaskBuilder(
	privategkemaster_contract.InputPrivateGKEMasterComponentNameFilterTaskID,
	priorityForPrivateGKEMasterGroup+1000,
	"Master component names",
).
	WithDefaultValueConstant([]string{"@any", "-kube-apiserver"}, true).
	WithAllowAddAll(false).
	WithAllowRemoveAll(false).
	WithAllowCustomValue(true).
	WithOptionsFunc(func(ctx context.Context, previousValues []string) ([]inspectionmetadata.SetParameterFormFieldOptionItem, error) {
		options := []inspectionmetadata.SetParameterFormFieldOptionItem{
			{ID: "@any", Description: "[Alias]Matches any component name"},
		}
		componentNames := []string{
			"cloud-controller-manager", "cluster-autoscaler", "clustermetrics", "containerd",
			"customer-logs-exporter", "directpath-agent", "directpath-agent-v2", "envoy",
			"gcp-controller-manager", "gke-audit-proxy", "gke-common-webhooks", "gke-kspan-proxy",
			"gke-kubestore-collector", "gke-master-cis-scanner", "gke-master-healthcheck",
			"gke-master-maintenance", "gke-master-monitoring", "gke-volume-populator-controller",
			"google-osconfig-agent", "helm-security", "konnectivity-server", "kube-addon-manager",
			"kube-apiserver", "kube-controller-manager", "kube-master-configuration",
			"kube-master-installation", "kube-scheduler", "kubelet", "l7-lb-controller",
			"l7-lb-controller-neg", "maintenance-controller", "managed-certificate-controller",
			"master-metrics-agent", "master-prom-to-sd-monitor", "masterfluentbit",
			"metrics-server-nanny", "pdcsi-controller",
		}

		for _, name := range componentNames {
			options = append(options, inspectionmetadata.SetParameterFormFieldOptionItem{
				ID:          name,
				Description: "Matches logs from " + name,
			})
		}
		return options, nil
	}).
	WithDescription("Component names to query (e.g. kube-apiserver, kube-scheduler...etc)").
	WithValidator(validatePrivateGKEMasterComponentFilter).
	WithConverter(convertPrivateGKEMasterComponentFilter).
	Build()

func validatePrivateGKEMasterComponentFilter(ctx context.Context, value []string) (string, error) {
	strFilter := strings.Join(value, " ")
	result, err := gcpqueryutil.ParseSetFilter(strFilter, map[string][]string{}, true, true, true)
	if err != nil {
		return "", err
	}
	return result.ValidationError, nil
}

func convertPrivateGKEMasterComponentFilter(ctx context.Context, value []string) (*gcpqueryutil.SetFilterParseResult, error) {
	strFilter := strings.Join(value, " ")
	result, err := gcpqueryutil.ParseSetFilter(strFilter, map[string][]string{}, true, true, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}
