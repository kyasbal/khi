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

package privatecsmcp_impl

import (
	"testing"

	form_task_test "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/formtask/test"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecsmcp_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecsmcp/contract"
)

func TestInputCSMTenantProjectIDTask(t *testing.T) {
	form_task_test.TestTextForms(t, "tenant project ID", InputCSMTenantProjectIDTask, []*form_task_test.TextFormTestCase{
		{
			Name:          "empty input",
			Input:         "",
			ExpectedValue: "",
			ExpectedFormField: inspectionmetadata.TextParameterFormField{
				ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
					ID:          "private-csmcp-input-csm-tenant-project-id",
					Type:        "Text",
					Label:       "CSM Tenant Project ID",
					Description: "The project ID of the CSM Tenant where Cloud Run metrics reside.",
					HintType:    inspectionmetadata.Error,
					Hint:        "Project ID must match `^[0-9a-z\\.:\\-]+$`",
				},
				Default:          "",
				ValidationTiming: inspectionmetadata.Blur,
			},
		},
	})
}

func TestInputCSMCPCloudRunServiceNameTask(t *testing.T) {
	mockAutocompleteSuccess := tasktest.StubTaskFromReferenceID(privatecsmcp_contract.AutocompleteCSMCPCloudRunServiceNameTaskID.Ref(), &inspectioncore_contract.AutocompleteResult[string]{
		Values: []string{"service-a", "service-b"},
		Error:  "",
	}, nil)

	mockAutocompleteError := tasktest.StubTaskFromReferenceID(privatecsmcp_contract.AutocompleteCSMCPCloudRunServiceNameTaskID.Ref(), &inspectioncore_contract.AutocompleteResult[string]{
		Values: []string{},
		Error:  "API error",
	}, nil)

	mockAutocompleteHint := tasktest.StubTaskFromReferenceID(privatecsmcp_contract.AutocompleteCSMCPCloudRunServiceNameTaskID.Ref(), &inspectioncore_contract.AutocompleteResult[string]{
		Values: []string{},
		Hint:   "Need project ID",
	}, nil)

	form_task_test.TestTextForms(t, "service name", InputCSMCPCloudRunServiceNameTask, []*form_task_test.TextFormTestCase{
		{
			Name:          "with autocomplete values, no previous value",
			Input:         "a",
			ExpectedValue: "a",
			Dependencies:  []coretask.UntypedTask{mockAutocompleteSuccess},
			ExpectedFormField: inspectionmetadata.TextParameterFormField{
				ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
					ID:          "private-csmcp-input-csm-cp-cloud-run-service-name",
					Type:        "Text",
					Label:       "CSM CP Cloud Run Service Name",
					Description: "The name of the Cloud Run service for CSM Control Plane.",
					HintType:    inspectionmetadata.None,
					Hint:        "",
				},
				Default:          "service-a",
				Suggestions:      []string{"service-a"},
				ValidationTiming: inspectionmetadata.Blur,
			},
		},
		{
			Name:          "with autocomplete error",
			Input:         "",
			ExpectedValue: "",
			Dependencies:  []coretask.UntypedTask{mockAutocompleteError},
			ExpectedFormField: inspectionmetadata.TextParameterFormField{
				ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
					ID:          "private-csmcp-input-csm-cp-cloud-run-service-name",
					Type:        "Text",
					Label:       "CSM CP Cloud Run Service Name",
					Description: "The name of the Cloud Run service for CSM Control Plane.",
					HintType:    inspectionmetadata.Error,
					Hint:        "Cloud Run Service Name must not be empty",
				},
				Default:          "",
				Suggestions:      nil,
				ValidationTiming: inspectionmetadata.Blur,
			},
		},
		{
			Name:          "with autocomplete hint",
			Input:         "",
			ExpectedValue: "",
			Dependencies:  []coretask.UntypedTask{mockAutocompleteHint},
			ExpectedFormField: inspectionmetadata.TextParameterFormField{
				ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
					ID:          "private-csmcp-input-csm-cp-cloud-run-service-name",
					Type:        "Text",
					Label:       "CSM CP Cloud Run Service Name",
					Description: "The name of the Cloud Run service for CSM Control Plane.",
					HintType:    inspectionmetadata.Error,
					Hint:        "Cloud Run Service Name must not be empty",
				},
				Default:          "",
				Suggestions:      nil,
				ValidationTiming: inspectionmetadata.Blur,
			},
		},
	})
}
