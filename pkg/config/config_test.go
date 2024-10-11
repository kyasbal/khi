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

	// case1
	keys, err := config.GetMergeKeys(model.KubernetesObjectOperation{
		PluralKind: "foo",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, but got %d", len(keys))
	}
	if keys["path1"] != "key1" {
		t.Errorf("Expected key1 for path1, but got %s", keys["path1"])
	}
	if keys["path2"] != "key2" {
		t.Errorf("Expected key2 for path2, but got %s", keys["path2"])
	}
	if keys["path3"] != "key3" {
		t.Errorf("Expected key3 for path3, but got %s", keys["path3"])
	}

	// case2
	keys, err = config.GetMergeKeys(model.KubernetesObjectOperation{
		PluralKind: "foo",
		Namespace:  "not-baz",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, but got %d", len(keys))
	}
	if keys["path1"] != "key1" {
		t.Errorf("Expected key1 for path1, but got %s", keys["path1"])
	}
	if keys["path3"] != "key3" {
		t.Errorf("Expected key3 for path3, but got %s", keys["path3"])
	}

	// case3
	keys, err = config.GetMergeKeys(model.KubernetesObjectOperation{
		PluralKind: "not-foo",
		Namespace:  "not-baz",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(keys) != 1 {
		t.Errorf("Expected 1 key, but got %d", len(keys))
	}
	if keys["path3"] != "key3" {
		t.Errorf("Expected key3 for path3, but got %s", keys["path3"])
	}
}
