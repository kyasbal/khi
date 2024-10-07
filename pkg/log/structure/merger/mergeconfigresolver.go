package merger

import (
	"fmt"
	"log/slog"
	"sync"
)

type MergeArrayStrategy string

const MergeStrategyMerge MergeArrayStrategy = "merge"
const MergeStrategyReplace MergeArrayStrategy = "replace"

var warnShownPath = sync.Map{}

type MergeConfigResolver struct {
	Parent          *MergeConfigResolver
	MergeStrategies map[string]MergeArrayStrategy
	MergeKeys       map[string]string
}

func (r *MergeConfigResolver) GetMergeArrayStrategy(fieldPath string) MergeArrayStrategy {
	if strategy, found := r.MergeStrategies[fieldPath]; found {
		return strategy
	} else {
		if r.Parent != nil {
			return r.Parent.GetMergeArrayStrategy(fieldPath)
		}
		_, found := warnShownPath.LoadOrStore(fieldPath, struct{}{})
		if !found {
			slog.Debug(fmt.Sprintf("Merge strategy for %s is not defined. Use replace strategy.", fieldPath))
		}
		return MergeStrategyReplace
	}
}

func (r *MergeConfigResolver) GetMergeKey(fieldPath string) (string, error) {
	if key, found := r.MergeKeys[fieldPath]; found {
		return key, nil
	} else {
		if r.Parent != nil {
			return r.Parent.GetMergeKey(fieldPath)
		}
		return "", fmt.Errorf("merge key for %s was not found", fieldPath)
	}
}
