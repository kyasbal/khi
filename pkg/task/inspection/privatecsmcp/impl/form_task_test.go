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
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	form_task_test "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/formtask/test"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api/iamtoken"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
	privatecsmcp_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecsmcp/contract"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestInputCSMTenantProjectIDTask(t *testing.T) {
	testCases := []struct {
		name              string
		prepareContext    func() context.Context
		hasInput          bool
		input             string
		expectedValue     string
		expectedFormField inspectionmetadata.TextParameterFormField
	}{
		{
			name: "auto-fills single registered tenant project and provides suggestions",
			prepareContext: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("csm-service-tp"), "token-csm")
				return khictx.WithValue(context.Background(), privatecommon_contract.APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			hasInput:      false,
			input:         "",
			expectedValue: "csm-service-tp",
			expectedFormField: inspectionmetadata.TextParameterFormField{
				ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
					ID:          "private-csmcp-input-csm-tenant-project-id",
					Type:        "Text",
					Label:       "CSM Tenant Project ID",
					Description: "The project ID of the CSM Tenant where Cloud Run metrics reside.",
					HintType:    inspectionmetadata.None,
					Hint:        "",
				},
				Default:          "csm-service-tp",
				Suggestions:      []string{"csm-service-tp"},
				ValidationTiming: inspectionmetadata.Blur,
			},
		},
		{
			name: "valid input with multiple suggestions",
			prepareContext: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("csm-1-tp"), "token-1")
				injector.SetTokenFor(googlecloud.Project("csm-2-tp"), "token-2")
				return khictx.WithValue(context.Background(), privatecommon_contract.APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			hasInput:      true,
			input:         "csm-1-tp",
			expectedValue: "csm-1-tp",
			expectedFormField: inspectionmetadata.TextParameterFormField{
				ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
					ID:          "private-csmcp-input-csm-tenant-project-id",
					Type:        "Text",
					Label:       "CSM Tenant Project ID",
					Description: "The project ID of the CSM Tenant where Cloud Run metrics reside.",
					HintType:    inspectionmetadata.None,
					Hint:        "",
				},
				Default:          "",
				Suggestions:      []string{"csm-1-tp", "csm-2-tp"},
				ValidationTiming: inspectionmetadata.Blur,
			},
		},
		{
			name:           "empty input without injector returns validation error",
			prepareContext: context.Background,
			hasInput:       true,
			input:          "",
			expectedValue:  "",
			expectedFormField: inspectionmetadata.TextParameterFormField{
				ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
					ID:          "private-csmcp-input-csm-tenant-project-id",
					Type:        "Text",
					Label:       "CSM Tenant Project ID",
					Description: "The project ID of the CSM Tenant where Cloud Run metrics reside.",
					HintType:    inspectionmetadata.Error,
					Hint:        "Project ID must match `^[0-9a-z\\.:\\-]+$`",
				},
				Default:          "",
				Suggestions:      nil,
				ValidationTiming: inspectionmetadata.Blur,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := inspectiontest.WithDefaultTestInspectionTaskContext(tc.prepareContext())
			inputMap := map[string]any{}
			if tc.hasInput {
				inputMap[InputCSMTenantProjectIDTask.ID().ReferenceIDString()] = tc.input
			}
			val, metadata, err := inspectiontest.RunInspectionTask(
				ctx,
				InputCSMTenantProjectIDTask,
				inspectioncore_contract.TaskModeDryRun,
				inputMap,
			)
			if err != nil {
				t.Fatalf("unexpected error running task: %v", err)
			}
			if diff := cmp.Diff(tc.expectedValue, val); diff != "" {
				t.Errorf("task result value mismatch (-want +got):\n%s", diff)
			}

			formFields, found := typedmap.Get(metadata, inspectionmetadata.FormFieldSetMetadataKey)
			if !found {
				t.Fatalf("form field metadata not found")
			}
			field := formFields.DangerouslyGetField(InputCSMTenantProjectIDTask.UntypedID().GetUntypedReference().String())
			if diff := cmp.Diff(tc.expectedFormField, field, cmpopts.IgnoreFields(inspectionmetadata.ParameterFormFieldBase{}, "Priority", "ID", "Type")); diff != "" {
				t.Errorf("form field mismatch (-want +got):\n%s", diff)
			}
		})
	}
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
