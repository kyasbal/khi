package k8s

import (
	"fmt"
	"log/slog"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/merger"
)

type MergeConfigRegistry struct {
	defaultResolver      *merger.MergeConfigResolver
	mergeConfigResolvers map[string]*merger.MergeConfigResolver
}

func (r *MergeConfigRegistry) Register(apiVersion string, kind string, childResolver *merger.MergeConfigResolver) {
	mapKey := fmt.Sprintf("%s-%s", apiVersion, kind)
	if _, found := r.mergeConfigResolvers[mapKey]; found {
		slog.Error(fmt.Sprintf("Merge config for apiVersion: %s, kind:%s is already registered", apiVersion, kind))
	}
	r.mergeConfigResolvers[mapKey] = childResolver
}

func (r *MergeConfigRegistry) Get(apiVersion string, kind string) *merger.MergeConfigResolver {
	mapKey := fmt.Sprintf("%s-%s", apiVersion, kind)
	if resolver, found := r.mergeConfigResolvers[mapKey]; found {
		return resolver
	} else {
		slog.Warn(fmt.Sprintf("Merge config for apiVersion: %s, kind:%s was not found", apiVersion, kind))
		return r.defaultResolver
	}
}
