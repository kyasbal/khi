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

package inspectionmetadata

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func fieldWithIdAndPriorityForTest(id string, priority int) TextParameterFormField {
	return TextParameterFormField{
		ParameterFormFieldBase: ParameterFormFieldBase{
			ID:       id,
			Priority: priority,
		},
	}
}

func TestFormFieldSetShouldSortOnAddingNewField(t *testing.T) {
	fsActual := NewFormFieldSetMetadata()
	fsActual.SetField(fieldWithIdAndPriorityForTest("foo", 1))
	fsActual.SetField(fieldWithIdAndPriorityForTest("bar", 3))
	fsActual.SetField(fieldWithIdAndPriorityForTest("qux", 2))
	fsActual.SetField(fieldWithIdAndPriorityForTest("aqux", 2))

	fsExpected := &FormFieldSetMetadata{
		fields: []ParameterFormField{
			fieldWithIdAndPriorityForTest("bar", 3),
			fieldWithIdAndPriorityForTest("aqux", 2),
			fieldWithIdAndPriorityForTest("qux", 2),
			fieldWithIdAndPriorityForTest("foo", 1),
		},
	}

	if diff := cmp.Diff(fsActual, fsExpected, cmp.AllowUnexported(FormFieldSetMetadata{}), cmpopts.IgnoreFields(FormFieldSetMetadata{}, "fieldsLock")); diff != "" {
		t.Errorf("FieldSet has fields in unexpected shape\n%v", diff)
	}
}

func TestGetParameterFormFieldBase(t *testing.T) {
	testCases := []struct {
		name  string
		input ParameterFormField
		want  ParameterFormFieldBase
	}{
		{
			name: "group field",
			input: GroupParameterFormField{
				ParameterFormFieldBase: ParameterFormFieldBase{
					ID:    "group-1",
					Label: "Group 1",
					Type:  Group,
				},
			},
			want: ParameterFormFieldBase{
				ID:    "group-1",
				Label: "Group 1",
				Type:  Group,
			},
		},
		{
			name: "text field",
			input: TextParameterFormField{
				ParameterFormFieldBase: ParameterFormFieldBase{
					ID:    "text-1",
					Label: "Text 1",
					Type:  Text,
				},
			},
			want: ParameterFormFieldBase{
				ID:    "text-1",
				Label: "Text 1",
				Type:  Text,
			},
		},
		{
			name: "set field",
			input: SetParameterFormField{
				ParameterFormFieldBase: ParameterFormFieldBase{
					ID:    "set-1",
					Label: "Set 1",
					Type:  Set,
				},
			},
			want: ParameterFormFieldBase{
				ID:    "set-1",
				Label: "Set 1",
				Type:  Set,
			},
		},
		{
			name: "file field",
			input: FileParameterFormField{
				ParameterFormFieldBase: ParameterFormFieldBase{
					ID:    "file-1",
					Label: "File 1",
					Type:  File,
				},
			},
			want: ParameterFormFieldBase{
				ID:    "file-1",
				Label: "File 1",
				Type:  File,
			},
		},
		{
			name: "checkbox field",
			input: CheckboxParameterFormField{
				ParameterFormFieldBase: ParameterFormFieldBase{
					ID:    "checkbox-1",
					Label: "Checkbox 1",
					Type:  Checkbox,
				},
				Default:  true,
				Readonly: false,
			},
			want: ParameterFormFieldBase{
				ID:    "checkbox-1",
				Label: "Checkbox 1",
				Type:  Checkbox,
			},
		},
		{
			name: "list field",
			input: ListParameterFormField{
				ParameterFormFieldBase: ParameterFormFieldBase{
					ID:    "list-1",
					Label: "List 1",
					Type:  List,
				},
				MinCount: 1,
				MaxCount: 5,
			},
			want: ParameterFormFieldBase{
				ID:    "list-1",
				Label: "List 1",
				Type:  List,
			},
		},
		{
			name:  "unknown field",
			input: struct{}{},
			want:  ParameterFormFieldBase{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetParameterFormFieldBase(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GetParameterFormFieldBase() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetFileFieldIDs(t *testing.T) {
	testCases := []struct {
		name   string
		fields []ParameterFormField
		want   []string
	}{
		{
			name: "top-level file fields",
			fields: []ParameterFormField{
				FileParameterFormField{
					ParameterFormFieldBase: ParameterFormFieldBase{
						ID:   "file-1",
						Type: File,
					},
				},
				TextParameterFormField{
					ParameterFormFieldBase: ParameterFormFieldBase{
						ID:   "text-1",
						Type: Text,
					},
				},
				FileParameterFormField{
					ParameterFormFieldBase: ParameterFormFieldBase{
						ID:   "file-2",
						Type: File,
					},
				},
			},
			want: []string{"file-1", "file-2"},
		},
		{
			name: "file fields nested in group field",
			fields: []ParameterFormField{
				GroupParameterFormField{
					ParameterFormFieldBase: ParameterFormFieldBase{
						ID:   "group-1",
						Type: Group,
					},
					Children: []ParameterFormField{
						FileParameterFormField{
							ParameterFormFieldBase: ParameterFormFieldBase{
								ID:   "nested-file-1",
								Type: File,
							},
						},
						TextParameterFormField{
							ParameterFormFieldBase: ParameterFormFieldBase{
								ID:   "nested-text-1",
								Type: Text,
							},
						},
					},
				},
			},
			want: []string{"nested-file-1"},
		},
		{
			name: "file fields nested in list field",
			fields: []ParameterFormField{
				ListParameterFormField{
					ParameterFormFieldBase: ParameterFormFieldBase{
						ID:   "list-1",
						Type: List,
					},
					Items: []ListParameterFormFieldItem{
						{
							Key: "0",
							Field: FileParameterFormField{
								ParameterFormFieldBase: ParameterFormFieldBase{
									ID:   "list-file-0",
									Type: File,
								},
							},
						},
						{
							Key: "1",
							Field: FileParameterFormField{
								ParameterFormFieldBase: ParameterFormFieldBase{
									ID:   "list-file-1",
									Type: File,
								},
							},
						},
					},
				},
			},
			want: []string{"list-file-0", "list-file-1"},
		},
		{
			name: "file fields nested in group inside list field",
			fields: []ParameterFormField{
				ListParameterFormField{
					ParameterFormFieldBase: ParameterFormFieldBase{
						ID:   "nodes",
						Type: List,
					},
					Items: []ListParameterFormFieldItem{
						{
							Key: "0",
							Field: GroupParameterFormField{
								ParameterFormFieldBase: ParameterFormFieldBase{
									ID:   "nodes/0",
									Type: Group,
								},
								Children: []ParameterFormField{
									FileParameterFormField{
										ParameterFormFieldBase: ParameterFormFieldBase{
											ID:   "nodes/0/kubelet",
											Type: File,
										},
									},
									FileParameterFormField{
										ParameterFormFieldBase: ParameterFormFieldBase{
											ID:   "nodes/0/containerd",
											Type: File,
										},
									},
								},
							},
						},
					},
				},
			},
			want: []string{"nodes/0/kubelet", "nodes/0/containerd"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := &FormFieldSetMetadata{
				fields: tc.fields,
			}
			got := fs.GetFileFieldIDs()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GetFileFieldIDs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
