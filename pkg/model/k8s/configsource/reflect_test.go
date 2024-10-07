package configsource

import (
	"fmt"
	"log/slog"
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/merger"
)

func TestFromResourceTypeReflection(t *testing.T) {
	type TestCaseMapField struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	type InlineFields struct {
		InlineField1       string             `json:"inlineField1"`
		InlineArrayMerge   []TestCaseMapField `json:"inlineArrayMerge" patchStrategy:"merge" patchMergeKey:"name"`
		InlineArrayReplace []TestCaseMapField `json:"inlineArrayReplace"`
	}
	type testStructSecondLayer struct {
		Name   string             `json:"name"`
		Values []TestCaseMapField `json:"values" patchStrategy:"merge" patchMergeKey:"name"`
	}
	type testStructSecondLayerReplace struct {
		Name   string             `json:"name"`
		Values []TestCaseMapField `json:"values"`
	}
	type testStruct struct {
		Scalar         int                            `json:"scalar"`
		MapType        map[string]string              `json:"mapType"`
		PrimitiveArray []string                       `json:"primitive,omitempty" patchStrategy:"merge"`
		MergeWithName  []TestCaseMapField             `json:"mergename,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
		Replace        []TestCaseMapField             `json:"replacearray,omitempty"`
		MergeMerge     []testStructSecondLayer        `json:"mergemerge,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
		ReplaceMerge   []testStructSecondLayer        `json:"replacemerge,omitempty"`
		MergeReplace   []testStructSecondLayerReplace `json:"mergereplace,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
		ReplaceReplace []testStructSecondLayerReplace `json:"replacereplace,omitempty"`
		Inline         InlineFields                   `json:",inline"`
		PointerType    *testStructSecondLayer         `json:"pointerType"`
	}
	type recursiveStruct struct {
		Name      string            `json:"name"`
		Recursive []recursiveStruct `json:"recursive,omitempty" patchStategy:"merge" patchMergeKey:"name"`
	}
	type fieldTestCase struct {
		path     string
		strategy merger.MergeArrayStrategy
		mergeKey string
	}
	testCase := []struct {
		name                        string
		resourceType                interface{}
		fieldTestCases              []fieldTestCase
		wantErrorOnGenerateResolver bool
	}{
		{
			name:         "core v1 Pod can be registered",
			resourceType: corev1.Pod{},
		},
		{
			name:         "simple",
			resourceType: testStruct{},
			fieldTestCases: []fieldTestCase{
				{
					path:     ".primitive",
					strategy: merger.MergeStrategyMerge,
					mergeKey: "",
				}, {
					path:     ".mergename",
					strategy: merger.MergeStrategyMerge,
					mergeKey: "name",
				}, {
					path:     ".replacearray",
					strategy: merger.MergeStrategyReplace,
					mergeKey: "",
				}, {
					path:     ".mergemerge",
					strategy: merger.MergeStrategyMerge,
					mergeKey: "name",
				},
				{
					path:     ".mergemerge[].values",
					strategy: merger.MergeStrategyMerge,
					mergeKey: "name",
				}, {
					path:     ".mergereplace",
					strategy: merger.MergeStrategyMerge,
					mergeKey: "name",
				},
				{
					path:     ".mergereplace[].values",
					strategy: merger.MergeStrategyReplace,
					mergeKey: "",
				},
				{
					path:     ".replacemerge",
					strategy: merger.MergeStrategyReplace,
					mergeKey: "",
				},
				{
					path:     ".replacemerge[].values",
					strategy: merger.MergeStrategyMerge,
					mergeKey: "name",
				},
				{
					path:     ".replacereplace",
					strategy: merger.MergeStrategyReplace,
					mergeKey: "",
				},
				{
					path:     ".replacereplace[].values",
					strategy: merger.MergeStrategyReplace,
					mergeKey: "",
				},
				{
					path:     ".replacereplace[].values",
					strategy: merger.MergeStrategyReplace,
					mergeKey: "",
				},
				{
					path:     ".replacereplace[].values",
					strategy: merger.MergeStrategyReplace,
					mergeKey: "",
				},
				{
					path:     ".inlineArrayMerge",
					strategy: merger.MergeStrategyMerge,
					mergeKey: "name",
				},
				{
					path:     ".inlineArrayReplace",
					strategy: merger.MergeStrategyReplace,
					mergeKey: "",
				},
				{
					path:     ".pointerType.values",
					strategy: merger.MergeStrategyMerge,
					mergeKey: "name",
				},
			},
		},
		// Errornous
		{
			name:                        "recursive struct",
			resourceType:                recursiveStruct{},
			wantErrorOnGenerateResolver: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			resolver, err := FromResourceTypeReflection(tc.resourceType)
			if tc.wantErrorOnGenerateResolver {
				if err == nil {
					t.Errorf("an error was expected but no error returned")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			for key, strategy := range resolver.MergeStrategies {
				if strategy == merger.MergeStrategyReplace {
					slog.Info(fmt.Sprintf("%s -> %s\n", key, strategy))
				} else {
					mergeKey, err := resolver.GetMergeKey(key)
					if err != nil {
						t.Fatal(err)
					}
					slog.Info(fmt.Sprintf("%s -> %s (%s)\n", key, strategy, mergeKey))
				}
			}
			for _, field := range tc.fieldTestCases {
				t.Run(field.path, func(t *testing.T) {
					t.Run("GetMergeArrayStrategy", func(t *testing.T) {
						strategy := resolver.GetMergeArrayStrategy(field.path)
						if strategy != field.strategy {
							t.Errorf("expected %s, actual %s", field.strategy, strategy)
						}
					})
					t.Run("GetMergeKey", func(t *testing.T) {
						mergeKey, err := resolver.GetMergeKey(field.path)
						if field.strategy == merger.MergeStrategyReplace {
							if err == nil {
								t.Errorf("GetMergeKey in the array field with replace merge policy should return an error but no error returned")
							}
						} else {
							if err != nil {
								t.Fatal(err)
							}
							if mergeKey != field.mergeKey {
								t.Errorf("expected %s, actual %s", field.mergeKey, mergeKey)
							}
						}
					})
				})
			}
		})
	}
}
