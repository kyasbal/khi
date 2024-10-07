package config

import (
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
)

type MergeKeyMapping struct {
	Key             string `yaml:"key"`
	FieldPath       string `yaml:"fieldPath"`
	APIVersion      string `yaml:"apiVersion"`
	Kind            string `yaml:"kind"`
	Namespace       string `yaml:"namespace"`
	Name            string `yaml:"name"`
	SubResourceName string `yaml:"subResourceName"`
}

type ConfigFile struct {
	MergeKeys []MergeKeyMapping `yaml:"mergeKeys"`
}

// Get primary key used for merging YAML configurations
func (config *ConfigFile) GetMergeKeys(k8sOps model.KubernetesObjectOperation) (map[string]string, error) {
	var result map[string]string = make(map[string]string)
	for _, mapping := range config.MergeKeys {
		if mapping.APIVersion != "" && mapping.APIVersion != k8sOps.APIVersion {
			continue
		}
		if mapping.Namespace != "" && mapping.Name != k8sOps.Namespace {
			continue
		}
		if mapping.Kind != "" && mapping.Kind != k8sOps.PluralKind {
			continue
		}
		if mapping.Name != "" && mapping.Name != k8sOps.Name {
			continue
		}
		if mapping.SubResourceName != "" && mapping.SubResourceName != k8sOps.SubResourceName {
			continue
		}
		v, exist := result[mapping.FieldPath]
		if exist && v != mapping.Key {
			return nil, fmt.Errorf("mergeKey conflict! %s", mapping.FieldPath)
		}
		result[mapping.FieldPath] = mapping.Key
	}
	return result, nil
}
