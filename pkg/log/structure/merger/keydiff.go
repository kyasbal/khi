package merger

type keyDiff struct {
	OnlyInPrev  []string
	OnlyInPatch []string
	Both        []string
}

func NewKeyDiff(prevKeys []string, patchKeys []string) *keyDiff {
	prevKeysInMap := map[string]struct{}{}
	for _, prevKey := range prevKeys {
		prevKeysInMap[prevKey] = struct{}{}
	}
	bothInMap := map[string]struct{}{}
	patchKeysInMap := map[string]struct{}{}
	for _, patchKey := range patchKeys {
		if _, found := prevKeysInMap[patchKey]; found {
			bothInMap[patchKey] = struct{}{}
			delete(prevKeysInMap, patchKey)
		} else {
			patchKeysInMap[patchKey] = struct{}{}
		}
	}
	result := keyDiff{
		OnlyInPrev:  make([]string, 0),
		OnlyInPatch: make([]string, 0),
		Both:        make([]string, 0),
	}
	// To merge map, the priorities of fields are:
	// 1. Fields defined in the prev data
	// 2. append new fields added by the patch
	for _, prevKey := range prevKeys {
		if _, found := bothInMap[prevKey]; found {
			result.Both = append(result.Both, prevKey)
		} else {
			result.OnlyInPrev = append(result.OnlyInPrev, prevKey)
		}
	}
	for _, patchKey := range patchKeys {
		if _, found := patchKeysInMap[patchKey]; found {
			result.OnlyInPatch = append(result.OnlyInPatch, patchKey)
		}
	}
	return &result
}

func NewKeyDiffForArrayMerge(prevKeys []string, patchKeys []string) *keyDiff {
	prevKeysInMap := map[string]struct{}{}
	for _, prevKey := range prevKeys {
		prevKeysInMap[prevKey] = struct{}{}
	}
	bothInMap := map[string]struct{}{}
	patchKeysInMap := map[string]struct{}{}
	for _, patchKey := range patchKeys {
		if _, found := prevKeysInMap[patchKey]; found {
			bothInMap[patchKey] = struct{}{}
			delete(prevKeysInMap, patchKey)
		} else {
			patchKeysInMap[patchKey] = struct{}{}
		}
	}
	result := keyDiff{
		OnlyInPrev:  make([]string, 0),
		OnlyInPatch: make([]string, 0),
		Both:        make([]string, 0),
	}
	// To merge map, the priorities of fields are:
	// 1. Fields defined in the prev data
	// 2. append new fields added by the patch
	for _, patchKey := range patchKeys {
		if _, found := bothInMap[patchKey]; found {
			result.Both = append(result.Both, patchKey)
		} else {
			result.OnlyInPatch = append(result.OnlyInPatch, patchKey)
		}
	}
	for _, prevKey := range prevKeys {
		if _, found := bothInMap[prevKey]; !found {
			result.OnlyInPrev = append(result.OnlyInPrev, prevKey)
		}
	}

	return &result
}
