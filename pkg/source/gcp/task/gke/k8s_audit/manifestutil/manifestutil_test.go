package manifestutil

import (
	"context"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/testlog"
	"github.com/google/go-cmp/cmp"
)

func TestParseDeletionStatus(t *testing.T) {
	testCases := []struct {
		name           string
		inputYaml      string
		inputOperation *model.KubernetesObjectOperation
		expectedStatus DeletionStatus
	}{
		{
			name: "deleted with grace period 0",
			inputYaml: `metadata:
  deletionGracePeriodSeconds: 0
`,
			inputOperation: &model.KubernetesObjectOperation{
				Verb: enum.RevisionVerbUpdate,
			},
			expectedStatus: DeletionStatusDeleted,
		},
		{
			name: "it regards deleting when deletionGracePeriodSeconds is non zero",
			inputYaml: `metadata:
  deletionGracePeriodSeconds: 1
`,
			inputOperation: &model.KubernetesObjectOperation{
				Verb: enum.RevisionVerbUpdate,
			},
			expectedStatus: DeletionStatusDeleting,
		},
		{
			name: "it regards deleting because deletionGracePeriodSeconds is larger than 0 even deletionTimestamp is included",
			inputYaml: `metadata:
  deletionGracePeriodSeconds: 1
  deletionTimestamp: 2024-01-01T00:00:00Z
`,
			inputOperation: &model.KubernetesObjectOperation{
				Verb: enum.RevisionVerbUpdate,
			},
			expectedStatus: DeletionStatusDeleting,
		},
		{
			name: "deleted because metadata.deletionTimestamp is included",
			inputYaml: `metadata:
  deletionTimestamp: 2024-01-01T00:00:00Z`,
			inputOperation: &model.KubernetesObjectOperation{
				Verb: enum.RevisionVerbPatch,
			},
			expectedStatus: DeletionStatusDeleted,
		},
		{
			name:           "deleted by verb",
			inputYaml:      "",
			inputOperation: &model.KubernetesObjectOperation{Verb: enum.RevisionVerbDelete},
			expectedStatus: DeletionStatusDeleted,
		},
		{
			name:           "non deleted",
			inputYaml:      "",
			inputOperation: &model.KubernetesObjectOperation{Verb: enum.RevisionVerbUpdate},
			expectedStatus: DeletionStatusNonDefined,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader, err := testlog.New(testlog.BaseYaml(tc.inputYaml)).BuildReader()
			if err != nil {
				t.Fatal(err)
			}
			status := ParseDeletionStatus(context.Background(), reader, tc.inputOperation)
			if diff := cmp.Diff(tc.expectedStatus, status); diff != "" {
				t.Errorf("returned creation time doesn't match with the expected result\n%s", diff)
			}
		})
	}
}

func TestParseCreationTime(t *testing.T) {
	defaultTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		name         string
		inputYaml    string
		expectedTime time.Time
	}{
		{
			name: "the simplest deleted case",
			inputYaml: `metadata:
  creationTimestamp: "2022-01-02T00:00:00Z"`,
			expectedTime: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		}, {
			name: "metadata without the creationTimestamp",
			inputYaml: `metadata:
  deletionTimestamp: 2022-01-01T00:00:00Z`,
			expectedTime: defaultTime,
		}, {
			name:         "without the metadata",
			inputYaml:    `spec: foo`,
			expectedTime: defaultTime,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader, err := testlog.New(testlog.BaseYaml(tc.inputYaml)).BuildReader()
			if err != nil {
				t.Fatal(err)
			}
			creationTimestamp := ParseCreationTime(reader, defaultTime)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.expectedTime, creationTimestamp); diff != "" {
				t.Errorf("returned creation time doesn't match with the expected result\n%s", diff)
			}
		})
	}
}
