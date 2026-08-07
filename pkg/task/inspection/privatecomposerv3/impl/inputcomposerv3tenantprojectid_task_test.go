package privatecomposerv3_impl

import (
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	inspectiontest "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/test"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	tasktest "github.com/GoogleCloudPlatform/khi/pkg/core/task/test"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api/iamtoken"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestInputComposerV3TenantProjectIdTask(t *testing.T) {
	mockProjectIDTask := tasktest.StubTaskFromReferenceID(googlecloudcommon_contract.InputProjectIdTaskID.Ref(), "composer-project", nil)
	mockJustificationTask := tasktest.StubTaskFromReferenceID(privatecommon_contract.JustificationFormTaskID.Ref(), "b/12345678", nil)

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
				injector.SetTokenFor(googlecloud.Project("composer-env-tp"), "sample-token")
				return khictx.WithValue(context.Background(), privatecommon_contract.APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			hasInput:      false,
			input:         "",
			expectedValue: "composer-env-tp",
			expectedFormField: inspectionmetadata.TextParameterFormField{
				ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
					ID:          "privatecomposerv3.khi.google.com/input-composer-v3-tenant-project-id",
					Type:        "Text",
					Label:       "Managed Airflow 3 Tenant Project ID",
					Description: "Type the tenant project ID for the Managed Airflow 3 environment. You can find the tenant ID from the tenant project section in Google Admin. The project ID must end with '-tp'.",
					HintType:    inspectionmetadata.None,
					Hint:        "",
				},
				Default:          "composer-env-tp",
				Suggestions:      []string{"composer-env-tp"},
				ValidationTiming: inspectionmetadata.Blur,
			},
		},
		{
			name: "valid input with matching IAM token",
			prepareContext: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("tenant-1-tp"), "token-1")
				injector.SetTokenFor(googlecloud.Project("tenant-2-tp"), "token-2")
				return khictx.WithValue(context.Background(), privatecommon_contract.APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			hasInput:      true,
			input:         "tenant-1-tp",
			expectedValue: "tenant-1-tp",
			expectedFormField: inspectionmetadata.TextParameterFormField{
				ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
					ID:          "privatecomposerv3.khi.google.com/input-composer-v3-tenant-project-id",
					Type:        "Text",
					Label:       "Managed Airflow 3 Tenant Project ID",
					Description: "Type the tenant project ID for the Managed Airflow 3 environment. You can find the tenant ID from the tenant project section in Google Admin. The project ID must end with '-tp'.",
					HintType:    inspectionmetadata.None,
					Hint:        "",
				},
				Default:          "",
				Suggestions:      []string{"tenant-1-tp", "tenant-2-tp"},
				ValidationTiming: inspectionmetadata.Blur,
			},
		},
		{
			name: "missing IAM token injector returns error",
			prepareContext: func() context.Context {
				return context.Background()
			},
			hasInput:      true,
			input:         "tenant-tp",
			expectedValue: "",
			expectedFormField: inspectionmetadata.TextParameterFormField{
				ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
					ID:          "privatecomposerv3.khi.google.com/input-composer-v3-tenant-project-id",
					Type:        "Text",
					Label:       "Managed Airflow 3 Tenant Project ID",
					Description: "Type the tenant project ID for the Managed Airflow 3 environment. You can find the tenant ID from the tenant project section in Google Admin. The project ID must end with '-tp'.",
					HintType:    inspectionmetadata.Error,
					Hint:        "IAMToken injector isn't set. Managed Airflow 3 parsers won't work unless you open KHI from Google Admin",
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
			dependencies := []coretask.UntypedTask{mockProjectIDTask, mockJustificationTask}
			inputMap := map[string]any{}
			if tc.hasInput {
				inputMap[InputComposerV3TenantProjectIdTask.ID().ReferenceIDString()] = tc.input
			}
			val, metadata, err := inspectiontest.RunInspectionTaskWithDependency(
				ctx,
				InputComposerV3TenantProjectIdTask,
				dependencies,
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
			field := formFields.DangerouslyGetField(InputComposerV3TenantProjectIdTask.UntypedID().GetUntypedReference().String())
			if diff := cmp.Diff(tc.expectedFormField, field, cmpopts.IgnoreFields(inspectionmetadata.ParameterFormFieldBase{}, "Priority", "ID", "Type")); diff != "" {
				t.Errorf("form field mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
