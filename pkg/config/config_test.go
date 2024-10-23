package config

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
)

func TestGetMergeKeys(t *testing.T) {
	config := ConfigFile{}
	config.MergeKeys = []MergeKeyMapping{
		{
			Key:             "key1",
			FieldPath:       "path1",
			APIVersion:      "",
			Kind:            "foo",
			Namespace:       "",
			Name:            "",
			SubResourceName: "",
		},
		{
			Key:             "key2",
			FieldPath:       "path2",
			APIVersion:      "",
			Kind:            "foo",
			Namespace:       "baz",
			Name:            "",
			SubResourceName: "",
		},
		{
			Key:             "key3",
			FieldPath:       "path3",
			APIVersion:      "",
			Kind:            "",
			Namespace:       "",
			Name:            "",
			SubResourceName: "",
		},
	}

	testCases := []struct {
		name     string
		k8sOps   model.KubernetesObjectOperation
		wantKeys map[string]string
		wantErr  bool
	}{
		{
			name: "matching all of the given merge keys",
			k8sOps: model.KubernetesObjectOperation{
				PluralKind: "foo",
			},
			wantKeys: map[string]string{
				"path1": "key1",
				"path2": "key2",
				"path3": "key3",
			},
			wantErr: false,
		},
		{
			name: "matching only the kind",
			k8sOps: model.KubernetesObjectOperation{
				PluralKind: "foo",
				Namespace:  "not-baz",
			},
			wantKeys: map[string]string{
				"path1": "key1",
				"path3": "key3",
			},
			wantErr: false,
		},
		{
			name: "not matching any selector but the field path",
			k8sOps: model.KubernetesObjectOperation{
				PluralKind: "not-foo",
				Namespace:  "not-baz",
			},
			wantKeys: map[string]string{
				"path3": "key3",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotKeys, err := config.GetMergeKeys(tc.k8sOps)
			if (err != nil) != tc.wantErr {
				t.Errorf("got %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(gotKeys) != len(tc.wantKeys) {
				t.Errorf("got %d keys, want %d", len(gotKeys), len(tc.wantKeys))
			}
			for k, v := range tc.wantKeys {
				if gotKeys[k] != v {
					t.Errorf("for %q, got %q, want %q", k, gotKeys[k], v)
				}
			}
		})
	}
}
