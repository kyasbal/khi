package v2commonlogparse

import (
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/k8saudittask"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/types"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/testlog"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/testtask"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseResourceSpecificParserInput(t *testing.T) {
	baseLog := `insertId: foo
protoPayload:
  authenticationInfo: 
    principalEmail: user@example.com
  methodName: io.k8s.core.v1.pods.create
  resourceName: core/v1/namespaces/default/pods/my-pod
  status:
    code: 200
timestamp: 2024-01-01T00:00:00+09:00`
	testCases := []struct {
		name            string
		testLogOpts     []testlog.TestLogOpt
		expectedResult  *types.ResourceSpecificParserInput
		expectedErr     bool
		expectedErrMsg  string
		skipSubResource bool
	}{
		{
			name:        "Parse basic k8s audit logs",
			testLogOpts: []testlog.TestLogOpt{},
			expectedResult: &types.ResourceSpecificParserInput{
				ResourceName:   "core/v1/namespaces/default/pods/my-pod",
				MethodName:     "io.k8s.core.v1.pods.create",
				PrincipalEmail: "user@example.com",
				Operation: &model.KubernetesObjectOperation{
					APIVersion: "core/v1",
					PluralKind: "pods",
					Namespace:  "default",
					Name:       "my-pod",
					Verb:       enum.RevisionVerbCreate,
				},
				Code: 200,
			},
			expectedErr:    false,
			expectedErrMsg: "",
		},
		{
			name: "Status subresource is merged with the parent",
			testLogOpts: []testlog.TestLogOpt{
				testlog.StringField("protoPayload.resourceName", "core/v1/namespaces/default/pods/my-pod/status"),
				testlog.StringField("protoPayload.methodName", "io.k8s.core.v1.pods.binding.patch"),
			},
			expectedResult: &types.ResourceSpecificParserInput{
				ResourceName:   "core/v1/namespaces/default/pods/my-pod/status",
				MethodName:     "io.k8s.core.v1.pods.binding.patch",
				PrincipalEmail: "user@example.com",
				Operation: &model.KubernetesObjectOperation{
					APIVersion: "core/v1",
					PluralKind: "pods",
					Namespace:  "default",
					Name:       "my-pod",
					Verb:       enum.RevisionVerbPatch,
				},
				Code: 200,
			},
			expectedErr:    false,
			expectedErrMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tl := testlog.New(testlog.BaseYaml(baseLog))
			tl = tl.With(tc.testLogOpts...)
			result, err := ParseResourceSpecificParserInputWithoutResourceBody(context.Background(), tl.MustBuildLogEntity(&log.UnreachableCommonFieldExtractor{}))

			if tc.expectedErr && err == nil {
				t.Errorf("expected an error returned, but got nil")
			} else {
				if diff := cmp.Diff(tc.expectedResult, result, cmpopts.IgnoreFields(types.ResourceSpecificParserInput{}, "Log")); diff != "" {
					t.Errorf("the result is not valid:\n%s", diff)
				}
				if err != nil {
					t.Errorf(err.Error())
				}
			}
		})
	}
}

func TestPrestepParseTaskFinishWithSuccess(t *testing.T) {
	baseLog := `insertId: foo
protoPayload:
  authenticationInfo: 
    principalEmail: user@example.com
  methodName: io.k8s.core.v1.pods.create
  resourceName: core/v1/namespaces/default/pods/my-pod
  status:
    code: 200
timestamp: 2024-01-01T00:00:00+09:00`
	tl := testlog.New(testlog.BaseYaml(baseLog))
	logs := []*log.LogEntity{}
	for i := 0; i < 100; i++ {
		logs = append(logs, tl.MustBuildLogEntity(&log.UnreachableCommonFieldExtractor{}))
	}
	expectedResult := &types.ResourceSpecificParserInput{
		ResourceName:   "core/v1/namespaces/default/pods/my-pod",
		MethodName:     "io.k8s.core.v1.pods.create",
		PrincipalEmail: "user@example.com",
		Operation: &model.KubernetesObjectOperation{
			APIVersion: "core/v1",
			PluralKind: "pods",
			Namespace:  "default",
			Name:       "my-pod",
			Verb:       enum.RevisionVerbCreate,
		},
		Code: 200,
	}

	res, err := testtask.RunSingleTask[[]*types.ResourceSpecificParserInput](Task, task.TaskModeRun,
		testtask.PriorTaskResultFromID(k8saudittask.K8sAuditQueryTaskId, logs))
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(res) != len(logs) {
		t.Errorf("result count mismatch")
	}
	for _, r := range res {
		if diff := cmp.Diff(expectedResult, r, cmpopts.IgnoreFields(types.ResourceSpecificParserInput{}, "Log")); diff != "" {
			t.Errorf("the result is not valid:\n%s", diff)
		}
	}
}

func TestPrestepParseIgnoreErrornousLogs(t *testing.T) {
	baseLog := `insertId: foo
protoPayload:
  authenticationInfo: 
    principalEmail: user@example.com
  resourceName: core/v1/namespaces/default/pods/my-pod
  status:
    code: 200
timestamp: 2024-01-01T00:00:00+09:00`
	tl := testlog.New(testlog.BaseYaml(baseLog))
	logs := []*log.LogEntity{}
	for i := 0; i < 100; i++ {
		tl2 := tl
		if i%2 == 0 {
			// only the even index log the valid log
			tl2 = tl2.With(testlog.StringField("protoPayload.methodName", "io.k8s.core.v1.pods.create"))
		}
		logs = append(logs, tl2.MustBuildLogEntity(&log.UnreachableCommonFieldExtractor{}))
	}
	expectedResult := &types.ResourceSpecificParserInput{
		ResourceName:   "core/v1/namespaces/default/pods/my-pod",
		MethodName:     "io.k8s.core.v1.pods.create",
		PrincipalEmail: "user@example.com",
		Operation: &model.KubernetesObjectOperation{
			APIVersion: "core/v1",
			PluralKind: "pods",
			Namespace:  "default",
			Name:       "my-pod",
			Verb:       enum.RevisionVerbCreate,
		},
		Code: 200,
	}

	res, err := testtask.RunSingleTask[[]*types.ResourceSpecificParserInput](Task, task.TaskModeRun,
		testtask.PriorTaskResultFromID(k8saudittask.K8sAuditQueryTaskId, logs))
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(res) != 50 {
		t.Errorf("result count mismatch")
	}
	for _, r := range res {
		if diff := cmp.Diff(expectedResult, r, cmpopts.IgnoreFields(types.ResourceSpecificParserInput{}, "Log")); diff != "" {
			t.Errorf("the result is not valid:\n%s", diff)
		}
	}

}
