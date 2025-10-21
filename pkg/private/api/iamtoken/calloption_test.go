package iamtoken

import (
	"context"
	"net/http"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestIamTokenCallOptionInjectorOption_ApplyToCallContext(t *testing.T) {
	const token = "test-token"
	option := &iamTokenCallOptionInjectorOption{token: token}
	ctx := context.Background()

	newCtx := option.ApplyToCallContext(ctx, nil)

	md, ok := metadata.FromOutgoingContext(newCtx)
	if !ok {
		t.Fatal("metadata not found in context")
	}

	tokens := md.Get(iamTokenKey)
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}

	if tokens[0] != token {
		t.Errorf("expected token %q, got %q", token, tokens[0])
	}
}

func TestIamTokenCallOptionInjectorOption_ApplyToRawHTTPHeader(t *testing.T) {
	const token = "test-token"
	option := &iamTokenCallOptionInjectorOption{token: token}
	header := http.Header{}

	option.ApplyToRawHTTPHeader(header, nil)

	if got := header.Get(iamTokenKey); got != token {
		t.Errorf("expected token %q, got %q", token, got)
	}
}
