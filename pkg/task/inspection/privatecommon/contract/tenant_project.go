// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
