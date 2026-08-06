package privatecommon_contract

import (
	"context"
	"slices"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common"
	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
)

// RegisteredTenantProjectIDs returns registered project IDs ending with "tp" from the IAM token injector in context.
func RegisteredTenantProjectIDs(ctx context.Context) []string {
	injector, err := khictx.GetValue(ctx, APIClientIAMTokenInjectorOptionContextKey)
	if err != nil || injector == nil {
		return nil
	}

	projectIDs := injector.RegisteredProjectIDs()
	var candidates []string
	for _, pid := range projectIDs {
		if strings.HasSuffix(pid, "tp") {
			candidates = append(candidates, pid)
		}
	}
	slices.Sort(candidates)
	return candidates
}

// TenantProjectIDSuggestionsProvider provides autocomplete suggestions for tenant project ID form fields.
// Retrieves registered projects ending with "tp" and sorts them against the current input value using Levenshtein distance.
func TenantProjectIDSuggestionsProvider(ctx context.Context, value string, previousValues []string) ([]string, error) {
	candidates := RegisteredTenantProjectIDs(ctx)
	if len(candidates) == 0 {
		return nil, nil
	}
	return common.SortForAutocomplete(value, candidates), nil
}

// TenantProjectIDDefaultValueProvider determines the default value for tenant project ID form fields.
// Returns the first element of previousValues if present, or the single registered tenant project ID if exactly one exists.
func TenantProjectIDDefaultValueProvider(ctx context.Context, previousValues []string) (string, error) {
	if len(previousValues) > 0 {
		return previousValues[0], nil
	}
	candidates := RegisteredTenantProjectIDs(ctx)
	if len(candidates) == 1 {
		return candidates[0], nil
	}
	return "", nil
}
