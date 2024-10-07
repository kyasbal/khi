package task

func newLocalCachedTaskRunnerForSingleTask(target Definition, cache TaskVariableCache, dependencies ...Definition) (*LocalRunner, error) {
	sourceDs, err := NewSet([]Definition{target})
	if err != nil {
		return nil, err
	}
	availableDs, err := NewSet(dependencies)
	if err != nil {
		return nil, err
	}

	resolved, err := sourceDs.ResolveTask(availableDs)
	if err != nil {
		return nil, err
	}

	localRunner, err := NewLocalRunner(resolved)
	if err != nil {
		return nil, err
	}
	localRunner.WithCacheProvider(cache)
	return localRunner, nil
}
