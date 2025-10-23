package iamtoken

import (
	"net/http"
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"google.golang.org/grpc/metadata"
)

func TestIamTokenCallOptionInjectorOption_ApplyToCallContext(t *testing.T) {
	const defaultToken = "test-default-iam-token"
	testCases := []struct {
		desc      string
		prepare   func() *IAMTokenCallOptionInjectorOption
		container googlecloud.ResourceContainer
		wantToken string
	}{
		{
			desc: "from default token",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				return New(defaultToken)
			},
			container: googlecloud.Project("foo"),
			wantToken: defaultToken,
		},
		{
			desc: "from container specific token",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				option := New(defaultToken)
				option.SetTokenFor(googlecloud.Project("foo"), "foo-token")
				return option
			},
			container: googlecloud.Project("foo"),
			wantToken: "foo-token",
		},
		{
			desc: "fallback to default token",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				option := New(defaultToken)
				option.SetTokenFor(googlecloud.Project("bar"), "bar-token")
				return option
			},
			container: googlecloud.Project("foo"),
			wantToken: defaultToken,
		},
		{
			desc: "nil container",
			prepare: func() *IAMTokenCallOptionInjectorOption {
				return New(defaultToken)
			},
			container: nil, // Default project just for getting locations...etc not different by projects
			wantToken: defaultToken,
		},
	}
	for _, tc := range testCases {
		option := tc.prepare()
		t.Run("ApplyToCallContext", func(t *testing.T) {
			t.Run(tc.desc, func(t *testing.T) {
				ctx := t.Context()
				ctx = option.ApplyToCallContext(ctx, tc.container)

				md, ok := metadata.FromOutgoingContext(ctx)
				if !ok {
					t.Fatal("metadata not found in context")
				}

				tokens := md.Get(iamTokenKey)
				if len(tokens) != 1 {
					t.Fatalf("expected 1 token, got %d", len(tokens))
				}

				if tokens[0] != tc.wantToken {
					t.Errorf("expected tc.wantToken %q, got %q", tc.wantToken, tokens[0])
				}
			})
		})

		t.Run("ApplyToRawHTTPHeader", func(t *testing.T) {
			t.Run(tc.desc, func(t *testing.T) {
				header := http.Header{}

				option.ApplyToRawHTTPHeader(header, tc.container)

				if got := header.Get(iamTokenKey); got != tc.wantToken {
					t.Errorf("expected token %q, got %q", tc.wantToken, got)
				}
			})
		})

	}
}
