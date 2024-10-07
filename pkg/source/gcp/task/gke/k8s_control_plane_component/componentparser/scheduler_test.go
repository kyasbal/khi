package componentparser

import (
	"context"
	"testing"

	log_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/log"
	"github.com/google/go-cmp/cmp"
)

func TestPodRelatedLogsToResourcePath(t *testing.T) {
	testCases := []struct {
		testName      string
		inputLog      string
		expectedPath  string
		expectedError bool
	}{
		{
			testName: "Normal case",
			inputLog: `insertId: zi5s4vf0ywd0ou9w
jsonPayload:
  message: '"Add event for scheduled pod" pod="1-2-deployment-update/nginx-deployment-non-surge-5455c7f485-2cfpz"'
  pid: "11"
logName: projects/khi-testing/logs/container.googleapis.com%2Fscheduler
receiveTimestamp: "2024-08-19T10:31:14.802511598Z"
resource:
  labels:
    cluster_name: gke-basic-1
    component_location: us-central1-a
    component_name: scheduler
    location: us-central1-a
    project_id: khi-testing
  type: k8s_control_plane_component
severity: INFO
sourceLocation:
  file: eventhandlers.go
  line: "197"
timestamp: "2024-08-19T10:31:12.865780Z"`,
			expectedPath:  "core/v1#pod#1-2-deployment-update#nginx-deployment-non-surge-5455c7f485-2cfpz",
			expectedError: false,
		},
		{
			testName: "Missing pod field",
			inputLog: `insertId: 1hyqhhwvyaqo49zm
jsonPayload:
  message: To require authentication configuration lookup to succeed, set --authentication-tolerate-lookup-failure=false
  pid: "11"
logName: projects/khi-testing/logs/container.googleapis.com%2Fscheduler
receiveTimestamp: "2024-08-19T10:07:36.527462487Z"
resource:
  labels:
    cluster_name: gke-basic-1
    component_location: us-central1-a
    component_name: scheduler
    location: us-central1-a
    project_id: khi-testing
  type: k8s_control_plane_component
severity: WARNING
sourceLocation:
  file: authentication.go
  line: "370"
timestamp: "2024-08-19T10:06:31.833958Z"`,
			expectedPath:  "",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			parser := &SchedulerComponentParser{}
			l := log_test.MustLogEntity(tc.inputLog)
			path, err := parser.podRelatedLogsToResourcePath(context.Background(), l)
			if tc.expectedError {
				if err == nil {
					t.Errorf("expected an error but no error returned")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
			if diff := cmp.Diff(tc.expectedPath, path); diff != "" {
				t.Errorf("the result path is not valid:\nInput:\n%v\nActual:\n%s\nExpected:\n%s", tc.inputLog, path, tc.expectedPath)
			}
		})
	}
}
