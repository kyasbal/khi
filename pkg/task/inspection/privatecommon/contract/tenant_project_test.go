package privatecommon_contract

import (
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api/iamtoken"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRegisteredTenantProjectIDs(t *testing.T) {
	testCases := []struct {
		name    string
		prepare func() context.Context
		want    []string
	}{
		{
			name: "returns only projects ending with tp",
			prepare: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("main-project"), "token-1")
				injector.SetTokenFor(googlecloud.Project("composer-env-tp"), "token-2")
				injector.SetTokenFor(googlecloud.Project("csm-service-tp"), "token-3")
				return khictx.WithValue(context.Background(), APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			want: []string{"composer-env-tp", "csm-service-tp"},
		},
		{
			name: "no tp projects registered",
			prepare: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("main-project"), "token-1")
				return khictx.WithValue(context.Background(), APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			want: nil,
		},
		{
			name: "no injector in context",
			prepare: func() context.Context {
				return context.Background()
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.prepare()
			got := RegisteredTenantProjectIDs(ctx)
			if diff := cmp.Diff(tc.want, got, cmpopts.SortSlices(func(a, b string) bool { return a < b }), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("RegisteredTenantProjectIDs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTenantProjectIDSuggestionsProvider(t *testing.T) {
	testCases := []struct {
		name           string
		prepare        func() context.Context
		value          string
		previousValues []string
		want           []string
	}{
		{
			name: "empty input returns all matching candidates",
			prepare: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("project-a-tp"), "token-a")
				injector.SetTokenFor(googlecloud.Project("project-b-tp"), "token-b")
				injector.SetTokenFor(googlecloud.Project("main-project"), "token-main")
				return khictx.WithValue(context.Background(), APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			value:          "",
			previousValues: nil,
			want:           []string{"project-b-tp", "project-a-tp"},
		},
		{
			name: "prefix input sorts matching candidate first",
			prepare: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("alpha-tp"), "token-a")
				injector.SetTokenFor(googlecloud.Project("beta-tp"), "token-b")
				return khictx.WithValue(context.Background(), APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			value:          "bet",
			previousValues: nil,
			want:           []string{"beta-tp", "alpha-tp"},
		},
		{
			name: "no injector in context returns nil",
			prepare: func() context.Context {
				return context.Background()
			},
			value:          "test",
			previousValues: nil,
			want:           nil,
		},
		{
			name: "no candidates returns nil",
			prepare: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("main-project"), "token-main")
				return khictx.WithValue(context.Background(), APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			value:          "",
			previousValues: nil,
			want:           nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.prepare()
			got, err := TenantProjectIDSuggestionsProvider(ctx, tc.value, tc.previousValues)
			if err != nil {
				t.Fatalf("TenantProjectIDSuggestionsProvider() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("TenantProjectIDSuggestionsProvider() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTenantProjectIDDefaultValueProvider(t *testing.T) {
	testCases := []struct {
		name           string
		prepare        func() context.Context
		previousValues []string
		want           string
	}{
		{
			name: "prefers previous value if present",
			prepare: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("single-project-tp"), "token-1")
				return khictx.WithValue(context.Background(), APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			previousValues: []string{"prev-project-tp"},
			want:           "prev-project-tp",
		},
		{
			name: "auto-fills if exactly one tp project is registered and no previous value",
			prepare: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("single-project-tp"), "token-1")
				injector.SetTokenFor(googlecloud.Project("non-tp-project"), "token-2")
				return khictx.WithValue(context.Background(), APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			previousValues: nil,
			want:           "single-project-tp",
		},
		{
			name: "returns empty if multiple tp projects are registered and no previous value",
			prepare: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("project-1-tp"), "token-1")
				injector.SetTokenFor(googlecloud.Project("project-2-tp"), "token-2")
				return khictx.WithValue(context.Background(), APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			previousValues: nil,
			want:           "",
		},
		{
			name: "returns empty if no tp projects registered",
			prepare: func() context.Context {
				injector := iamtoken.NewInjector()
				injector.SetTokenFor(googlecloud.Project("non-tp-project"), "token-2")
				return khictx.WithValue(context.Background(), APIClientIAMTokenInjectorOptionContextKey, injector)
			},
			previousValues: nil,
			want:           "",
		},
		{
			name: "returns empty if no injector in context",
			prepare: func() context.Context {
				return context.Background()
			},
			previousValues: nil,
			want:           "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.prepare()
			got, err := TenantProjectIDDefaultValueProvider(ctx, tc.previousValues)
			if err != nil {
				t.Fatalf("TenantProjectIDDefaultValueProvider() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("TenantProjectIDDefaultValueProvider() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
