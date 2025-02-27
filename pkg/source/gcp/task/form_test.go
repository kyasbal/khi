// Copyright 2024 Google LLC
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

package task

import (
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/khi/pkg/common"
	"github.com/GoogleCloudPlatform/khi/pkg/inspection/metadata"
	"github.com/GoogleCloudPlatform/khi/pkg/inspection/metadata/form"
	inspection_task "github.com/GoogleCloudPlatform/khi/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/khi/pkg/parameters"
	"github.com/GoogleCloudPlatform/khi/pkg/source/gcp/query/queryutil"
	"github.com/GoogleCloudPlatform/khi/pkg/task"
	"github.com/google/go-cmp/cmp/cmpopts"

	form_test "github.com/GoogleCloudPlatform/khi/pkg/testutil/form"
	task_test "github.com/GoogleCloudPlatform/khi/pkg/testutil/task"
	"github.com/GoogleCloudPlatform/khi/pkg/testutil/testtask"
)

var testClusterNamePrefix = task_test.MockProcessorTaskFromTaskID(ClusterNamePrefixTaskID, "")

func TestProjectIdInput(t *testing.T) {
	form_test.TestTextForms(t, "gcp-project-id", InputProjectIdTask, []*form_test.FormTestCase{
		{
			Name:          "With valid project ID",
			Input:         "foo-project",
			ExpectedValue: "foo-project",
			Dependencies: []task.Definition{
				testClusterNamePrefix,
			},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					ID:          GCPPrefix + "input/project-id",
					Type:        form.Text,
					Label:       "Project ID",
					Description: "A project ID containing the cluster to inspect",
					HintType:    form.None,
				},
			},
		},
		{
			Name:          "With fixed project ID from environment variable",
			Input:         "foo-project",
			ExpectedValue: "bar-project",
			Dependencies: []task.Definition{
				testClusterNamePrefix,
			},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					ID:          GCPPrefix + "input/project-id",
					Type:        form.Text,
					Label:       "Project ID",
					Description: "A project ID containing the cluster to inspect",
					HintType:    form.None,
				},
				Readonly: true,
				Default:  "bar-project",
			},
			Before: func() {
				expectedFixedProjectId := "bar-project"
				parameters.Auth.FixedProjectID = &expectedFixedProjectId
			},
			After: func() {
				parameters.Auth.FixedProjectID = nil
			},
		},
		{
			Name:          "With invalid project ID",
			Input:         "A invalid project ID",
			ExpectedValue: "",
			Dependencies: []task.Definition{
				testClusterNamePrefix,
			},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					ID:          GCPPrefix + "input/project-id",
					Type:        form.Text,
					Label:       "Project ID",
					Description: "A project ID containing the cluster to inspect",
					HintType:    form.Error,
					Hint:        "Project ID must match `^*[0-9a-z\\.:\\-]+$`",
				},
			},
		},
		{
			Name:          "Spaces around project ID must be trimmed",
			Input:         "  project-foo   ",
			ExpectedValue: "project-foo",
			Dependencies: []task.Definition{
				testClusterNamePrefix,
			},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					ID:          GCPPrefix + "input/project-id",
					Type:        form.Text,
					Label:       "Project ID",
					Description: "A project ID containing the cluster to inspect",
					HintType:    form.None,
				},
			},
		},
		{
			Name:          "With valid old style project ID",
			Input:         "  deprecated.com:but-still-usable-project-id   ",
			ExpectedValue: "deprecated.com:but-still-usable-project-id",
			Dependencies: []task.Definition{
				testClusterNamePrefix,
			},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					ID:          GCPPrefix + "input/project-id",
					Description: "A project ID containing the cluster to inspect",
					Type:        "Text",
					Label:       "Project ID",
					HintType:    form.None,
				},
			},
		},
	})
}

func TestClusterNameInput(t *testing.T) {
	mockClusterNamesTask1 := task_test.MockProcessorTaskFromTaskID(AutocompleteClusterNamesTaskID, &AutocompleteClusterNameList{
		ClusterNames: []string{"foo-cluster", "bar-cluster"},
		Error:        "",
	})
	form_test.TestTextForms(t, "cluster name", InputClusterNameTask, []*form_test.FormTestCase{
		{
			Name:          "with valid cluster name",
			Input:         "foo-cluster",
			ExpectedValue: "foo-cluster",
			Dependencies:  []task.Definition{mockClusterNamesTask1, testClusterNamePrefix},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					ID:       GCPPrefix + "input/cluster-name",
					Type:     "Text",
					Label:    "Cluster name",
					HintType: form.None,
				},
				Suggestions: []string{"foo-cluster", "bar-cluster"},
				Default:     "foo-cluster",
			},
		},
		{
			Name:          "spaces around cluster name must be trimmed",
			Input:         "  foo-cluster   ",
			ExpectedValue: "foo-cluster",
			Dependencies:  []task.Definition{mockClusterNamesTask1, testClusterNamePrefix},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					ID:       GCPPrefix + "input/cluster-name",
					Type:     "Text",
					Label:    "Cluster name",
					HintType: form.None,
				},
				Suggestions: []string{"foo-cluster", "bar-cluster"},
				Default:     "foo-cluster",
			},
		},
		{
			Name:          "invalid cluster name",
			Input:         "An invalid cluster name",
			ExpectedValue: "foo-cluster",
			Dependencies:  []task.Definition{mockClusterNamesTask1, testClusterNamePrefix},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					ID:       GCPPrefix + "input/cluster-name",
					Type:     "Text",
					Label:    "Cluster name",
					HintType: form.Error,
					Hint:     "Cluster name must match `^[0-9a-z:\\-]+$`",
				},
				Suggestions: common.SortForAutocomplete("An invalid cluster name", []string{"foo-cluster", "bar-cluster"}),
				Default:     "foo-cluster",
			},
		},
		{
			Name:          "non existing cluster should show a hint",
			Input:         "nonexisting-cluster",
			ExpectedValue: "nonexisting-cluster",
			Dependencies:  []task.Definition{mockClusterNamesTask1, testClusterNamePrefix},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					ID:       GCPPrefix + "input/cluster-name",
					Type:     "Text",
					Label:    "Cluster name",
					Hint:     "Cluster `nonexisting-cluster` was not found in the specified project at this time. It works for the clusters existed in the past but make sure the cluster name is right if you believe the cluster should be there.",
					HintType: form.Warning,
				},
				Suggestions: []string{"foo-cluster", "bar-cluster"},
				Default:     "foo-cluster",
			},
		},
	})
}

func TestDurationInput(t *testing.T) {
	expectedDescription := ""
	expectedLabel := "Duration"
	expectedSuggestions := []string{"1m", "10m", "1h", "3h", "12h", "24h"}
	timezoneTaskUTC := task_test.MockProcessorTaskFromTaskID(TimeZoneShiftInputTaskID, time.UTC)
	timezoneTaskJST := task_test.MockProcessorTaskFromTaskID(TimeZoneShiftInputTaskID, time.FixedZone("", 9*3600))
	currentTimeTask1 := task_test.MockProcessorTaskFromTaskID(inspection_task.InspectionTimeProducer.ID().String(), time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC))
	endTimeTask := task_test.MockProcessorTaskFromTaskID(InputEndTimeTask.ID().String(), time.Date(2023, time.April, 1, 12, 0, 0, 0, time.UTC))

	form_test.TestTextForms(t, "duration", InputDurationTask, []*form_test.FormTestCase{
		{
			Name:          "With valid time duration",
			Input:         "10m",
			ExpectedValue: time.Duration(time.Minute) * 10,
			Dependencies:  []task.Definition{endTimeTask, currentTimeTask1, timezoneTaskUTC},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					HintType:    form.Info,
					Hint: `Query range:
2023-04-01T11:50:00Z ~ 2023-04-01T12:00:00Z
(UTC: 2023-04-01T11:50:00 ~ 2023-04-01T12:00:00)
(PDT: 2023-04-01T04:50:00 ~ 2023-04-01T05:00:00)`,
				},
				Suggestions: expectedSuggestions,
				Default:     "1h",
			},
		},
		{
			Name:          "With invalid time duration",
			Input:         "foo",
			ExpectedValue: time.Hour,
			Dependencies:  []task.Definition{endTimeTask, currentTimeTask1, timezoneTaskUTC},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					Hint:        "time: invalid duration \"foo\"",
					HintType:    form.Error,
				},
				Default:     "1h",
				Suggestions: expectedSuggestions,
			},
		},
		{
			Name:          "With invalid time duration(negative)",
			Input:         "-10m",
			ExpectedValue: time.Hour,
			Dependencies:  []task.Definition{endTimeTask, currentTimeTask1, timezoneTaskUTC},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					Hint:        "duration must be positive",
					HintType:    form.Error,
				},
				Suggestions: expectedSuggestions,
				Default:     "1h",
			},
		},
		{
			Name:          "with longer duration starting before than 30 days",
			Input:         "672h", // starting time will be 30 days before the inspection time
			ExpectedValue: time.Hour * 672,
			Dependencies:  []task.Definition{endTimeTask, currentTimeTask1, timezoneTaskUTC},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Type:        "Text",
					Label:       expectedLabel,
					Description: expectedDescription,
					Hint: `Specified time range starts from over than 30 days ago, maybe some logs are missing and the generated result could be incomplete.
This duration can be too long for big clusters and lead OOM. Please retry with shorter duration when your machine crashed.
Query range:
2023-03-04T12:00:00Z ~ 2023-04-01T12:00:00Z
(UTC: 2023-03-04T12:00:00 ~ 2023-04-01T12:00:00)
(PDT: 2023-03-04T05:00:00 ~ 2023-04-01T05:00:00)`,
					HintType: form.Info,
				},
				Suggestions: expectedSuggestions,
				Default:     "1h",
			},
		},
		{
			Name:          "With non UTC timezone",
			Input:         "1h",
			ExpectedValue: time.Hour,
			Dependencies:  []task.Definition{endTimeTask, currentTimeTask1, timezoneTaskJST},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Type:        "Text",
					Label:       expectedLabel,
					Description: expectedDescription,
					Hint: `Query range:
2023-04-01T20:00:00+09:00 ~ 2023-04-01T21:00:00+09:00
(UTC: 2023-04-01T11:00:00 ~ 2023-04-01T12:00:00)
(PDT: 2023-04-01T04:00:00 ~ 2023-04-01T05:00:00)`,
					HintType: form.Info,
				},
				Suggestions: expectedSuggestions,
				Default:     "1h",
			},
		},
	})
}

func TestInputEndtime(t *testing.T) {
	expectedDescription := "The endtime of query. Please input it in the format of RFC3339\n(example: 2006-01-02T15:04:05-07:00)"
	expectedLabel := "End time"
	expectedValue1, err := time.Parse(time.RFC3339, "2020-01-02T03:04:05Z")
	if err != nil {
		t.Errorf("unexpected error\n%s", err)
	}
	expectedValue2, err := time.Parse(time.RFC3339, "2020-01-02T00:00:00Z")
	timezoneTaskUTC := task_test.MockProcessorTaskFromTaskID(TimeZoneShiftInputTaskID, time.UTC)
	timezoneTaskJST := task_test.MockProcessorTaskFromTaskID(TimeZoneShiftInputTaskID, time.FixedZone("", 9*3600))

	if err != nil {
		t.Errorf("unexpected error\n%s", err)
	}
	form_test.TestTextForms(t, "endtime", InputEndTimeTask, []*form_test.FormTestCase{
		{
			Name:          "with empty",
			Input:         "",
			ExpectedValue: expectedValue1,
			Dependencies:  []task.Definition{inspection_task.TestInspectionTimeTaskProducer("2020-01-02T03:04:05Z"), timezoneTaskUTC},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					Hint:        "invalid time format. Please specify in the format of `2006-01-02T15:04:05-07:00`(RFC3339)",
					HintType:    form.Error,
				},
				Default:     "2020-01-02T03:04:05Z",
				Suggestions: []string{},
			},
		},
		{
			Name:          "with valid timestamp and UTC timezone",
			Input:         "2020-01-02T00:00:00Z",
			ExpectedValue: expectedValue2,
			Dependencies:  []task.Definition{inspection_task.TestInspectionTimeTaskProducer("2020-01-02T03:04:05Z"), timezoneTaskUTC},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					HintType:    form.None,
				},
				Suggestions: []string{},
				Default:     "2020-01-02T03:04:05Z",
			},
		},
		{
			Name:          "with valid timestamp and non UTC timezone",
			Input:         "2020-01-02T00:00:00Z",
			ExpectedValue: expectedValue2,
			Dependencies:  []task.Definition{inspection_task.TestInspectionTimeTaskProducer("2020-01-02T03:04:05Z"), timezoneTaskJST},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					HintType:    form.None,
				},
				Suggestions: []string{},
				Default:     "2020-01-02T12:04:05+09:00",
			},
		},
	})
}

func TestInputStartTime(t *testing.T) {
	duration, err := time.ParseDuration("1h30m")
	if err != nil {
		t.Fatal(err)
	}
	endTime, err := time.Parse(time.RFC3339, "2023-01-02T15:45:00Z")
	if err != nil {
		t.Fatal(err)
	}
	startTime, err := testtask.RunSingleTask[time.Time](InputStartTimeTask, inspection_task.TaskModeDryRun,
		testtask.PriorTaskResultFromID(inspection_task.MetadataVariableName, metadata.NewSet()),
		testtask.PriorTaskResultFromID(InputDurationTaskID, duration),
		testtask.PriorTaskResultFromID(InputEndTimeTaskID, endTime),
		testtask.PriorTaskResultFromID(TimeZoneShiftInputTaskID, time.UTC),
	)
	if err != nil {
		t.Errorf("unexpected error\n%v", err)
	}
	expectedTime, err := time.Parse(time.RFC3339, "2023-01-02T14:15:00Z")
	if err != nil {
		t.Errorf("unexpected error\n%v", err)
	}

	if startTime.String() != expectedTime.String() {
		t.Errorf("returned time is not matching with the expected value\n%s", startTime)
	}
}

func TestInputKindName(t *testing.T) {
	expectedDescription := ""
	expectedLabel := "Kind"
	form_test.TestTextForms(t, "kind", InputKindFilterTask, []*form_test.FormTestCase{
		{
			Input: "",
			ExpectedValue: &queryutil.SetFilterParseResult{
				Additives:       inputKindNameAliasMap["default"],
				Subtractives:    []string{},
				ValidationError: "",
				SubtractMode:    false,
			},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					HintType:    form.Error,
					Hint:        "kind filter can't be empty",
				},
				Readonly: false,
				Default:  "@default",
			},
		},
		{
			Input: "pods replicasets",
			ExpectedValue: &queryutil.SetFilterParseResult{
				Additives:       []string{"pods", "replicasets"},
				Subtractives:    []string{},
				ValidationError: "",
				SubtractMode:    false,
			},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					HintType:    form.None,
				},

				Readonly: false,
				Default:  "@default",
			},
		},
		{
			Input: "@invalid_alias",
			ExpectedValue: &queryutil.SetFilterParseResult{
				Additives:       inputKindNameAliasMap["default"],
				Subtractives:    []string{},
				ValidationError: "",
				SubtractMode:    false,
			}, ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					Hint:        "alias `invalid_alias` was not found",
					HintType:    form.Error,
				},
				Default:  "@default",
				Readonly: false,
			},
		},
	}, cmpopts.SortSlices(func(a string, b string) bool {
		return strings.Compare(a, b) > 0
	}))
}

func TestInputNamespaces(t *testing.T) {
	expectedDescription := ""
	expectedLabel := "Namespaces"
	form_test.TestTextForms(t, "namespaces", InputNamespaceFilterTask, []*form_test.FormTestCase{
		{
			Input: "",
			ExpectedValue: &queryutil.SetFilterParseResult{
				Additives: []string{
					"#namespaced",
					"#cluster-scoped",
				},
				Subtractives:    []string{},
				ValidationError: "",
				SubtractMode:    false,
			},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					HintType:    form.Error,
					Hint:        "namespace filter can't be empty",
				},
				Readonly: false,
				Default:  "@all_cluster_scoped @all_namespaced",
			},
		},
		{
			Input: "kube-system default",
			ExpectedValue: &queryutil.SetFilterParseResult{
				Additives:       []string{"kube-system", "default"},
				Subtractives:    []string{},
				ValidationError: "",
				SubtractMode:    false,
			},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					HintType:    form.None,
				},
				Readonly: false,
				Default:  "@all_cluster_scoped @all_namespaced",
			},
		},
		{
			Input: "@all_cluster_scoped @all_namespaced",
			ExpectedValue: &queryutil.SetFilterParseResult{
				Additives:       []string{"#namespaced", "#cluster-scoped"},
				Subtractives:    []string{},
				ValidationError: "",
				SubtractMode:    false,
			}, ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       expectedLabel,
					Description: expectedDescription,
					HintType:    form.None,
				},
				Readonly: false,
				Default:  "@all_cluster_scoped @all_namespaced",
			},
		},
	}, cmpopts.SortSlices(func(a string, b string) bool {
		return strings.Compare(a, b) > 0
	}))
}

func TestNodeNameFiltertask(t *testing.T) {
	wantLabelName := "Node names"
	wantDescription := "A space-separated list of node name substrings used to collect node-related logs. If left blank, KHI gathers logs from all nodes in the cluster."
	form_test.TestTextForms(t, "node-name", InputNodeNameFilterTask, []*form_test.FormTestCase{
		{
			Name:          "With an empty input",
			Input:         "",
			ExpectedValue: []string{},
			Dependencies:  []task.Definition{},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       wantLabelName,
					Description: wantDescription,
					HintType:    form.None,
				},
				Readonly: false,
			},
		},
		{
			Name:          "With a single node name substring",
			Input:         "node-name-1",
			ExpectedValue: []string{"node-name-1"},
			Dependencies:  []task.Definition{},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       wantLabelName,
					Description: wantDescription,
					HintType:    form.None,
				},
			},
		},
		{
			Name:          "With multiple node name substrings",
			Input:         "node-name-1 node-name-2 node-name-3",
			ExpectedValue: []string{"node-name-1", "node-name-2", "node-name-3"},
			Dependencies:  []task.Definition{},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       wantLabelName,
					Description: wantDescription,
					HintType:    form.None,
				},
			},
		},
		{
			Name:          "With invalid node name substring",
			Input:         "node-name-1 invalid=node=name node-name-3",
			ExpectedValue: []string{},
			Dependencies:  []task.Definition{},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       wantLabelName,
					Description: wantDescription,
					Hint:        "substring `invalid=node=name` is not valid as a substring of node name",
					HintType:    form.Error,
				},
			},
		},
		{
			Name:          "With spaces around node name substring",
			Input:         "  node-name-1  node-name-2  ",
			ExpectedValue: []string{"node-name-1", "node-name-2"},
			Dependencies:  []task.Definition{},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					Label:       wantLabelName,
					Description: wantDescription,
					HintType:    form.None,
				},
			},
		},
	})
}

func TestLocationInput(t *testing.T) {
	form_test.TestTextForms(t, "gcp-location", InputLocationsTask, []*form_test.FormTestCase{
		{
			Name:          "With valid location",
			Input:         "asia-northeast1",
			ExpectedValue: "asia-northeast1",
			Dependencies:  []task.Definition{},
			ExpectedFormField: form.TextParameterFormField{
				ParameterFormFieldBase: form.ParameterFormFieldBase{
					ID:          GCPPrefix + "input/location",
					Type:        "Text",
					Label:       "Location",
					Description: "A location(regions) containing the environments to inspect",
					HintType:    form.None,
				},
				Readonly: false,
			},
		},
	})
}
