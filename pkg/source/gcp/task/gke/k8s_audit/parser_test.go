package k8s_audit

import (
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/adapter"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
	"github.com/google/go-cmp/cmp"
)

func TestIsDeletedResource(t *testing.T) {
	readerFactory := structure.NewReaderFactory(&structuredatastore.OnMemoryStructureDataStore{})
	testCases := []struct {
		name      string
		inputYaml string
		isDeleted bool
		wantError bool
	}{{
		name: "the simplest deleted case",
		inputYaml: `metadata:
  deletionTimestamp: 2022-01-01T00:00:00Z
`,
		isDeleted: true,
		wantError: false,
	}, {
		name: "the simplest non deleted case",
		inputYaml: `metadata:
  creationTimestamp: 2022-01-01T00:00:00Z
`,
		isDeleted: false,
		wantError: false,
	}, {
		name:      "without the metadata",
		inputYaml: `spec: foo`,
		isDeleted: false,
		wantError: true,
	},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			reader, err := readerFactory.NewReader(adapter.Yaml(tc.inputYaml))
			if err != nil {
				t.Fatal(err)
			}
			deleted, err := isDeletedResource(reader)
			if tc.wantError {
				if err == nil {
					t.Error("wants an error but no error returned")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if deleted != tc.isDeleted {
					t.Errorf("expected: %v, actual: %v", tc.isDeleted, deleted)
				}
			}
		})
	}
}

func TestGetCreationTimeFromManifest(t *testing.T) {
	readerFactory := structure.NewReaderFactory(&structuredatastore.OnMemoryStructureDataStore{})

	testCases := []struct {
		name         string
		inputYaml    string
		expectedTime time.Time
		wantError    bool
	}{
		{
			name: "the simplest deleted case",
			inputYaml: `metadata:
  creationTimestamp: "2022-01-01T00:00:00Z"`,
			expectedTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			wantError:    false,
		}, {
			name: "metadata without the creationTimestamp",
			inputYaml: `metadata:
  deletionTimestamp: 2022-01-01T00:00:00Z`,
			expectedTime: time.Time{},
			wantError:    true,
		}, {
			name:         "without the metadata",
			inputYaml:    `spec: foo`,
			expectedTime: time.Time{},
			wantError:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader, err := readerFactory.NewReader(adapter.Yaml(tc.inputYaml))
			if err != nil {
				t.Fatal(err)
			}
			creationTimestamp, err := getCreationTimeFromManifest(reader)
			if tc.wantError {
				if err == nil {
					t.Error("wants an error but no error returned")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if diff := cmp.Diff(tc.expectedTime, creationTimestamp); diff != "" {
					t.Errorf("returned creation time doesn't match with the expected result\n%s", diff)
				}
			}
		})
	}
}
