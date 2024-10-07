package config

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
	"github.com/stretchr/testify/assert"
)

var testConfigFolder = "../../test/model/"

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
	assert.Nil(t, err)
	assert.Equal(t, 3, len(keys))
	assert.Equal(t, "key1", keys["path1"])
	assert.Equal(t, "key2", keys["path2"])
	assert.Equal(t, "key3", keys["path3"])

	// case2
	keys, err = config.GetMergeKeys(model.KubernetesObjectOperation{
		PluralKind: "foo",
		Namespace:  "not-baz",
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(keys))
	assert.Equal(t, "key1", keys["path1"])
	assert.Equal(t, "key3", keys["path3"])

	// case3
	keys, err = config.GetMergeKeys(model.KubernetesObjectOperation{
		PluralKind: "not-foo",
		Namespace:  "not-baz",
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, "key3", keys["path3"])
}
