package iamtoken

import (
	"context"
	"net/http"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"google.golang.org/grpc/metadata"
)

const iamTokenKey = "x-goog-iam-authorization-token"

// iamTokenCallOptionInjectorOption is a CallOptionInjectorOption that injects an IAM token into API calls.
type iamTokenCallOptionInjectorOption struct {
	token string
}

func NewCallOptionInjectorOption(iamToken string) googlecloud.CallOptionInjectorOption {
	return &iamTokenCallOptionInjectorOption{iamToken}
}

// ApplyToCallContext implements googlecloud.CallOptionInjectorOption.
// It appends the IAM token to the outgoing context metadata for gRPC calls.
func (i *iamTokenCallOptionInjectorOption) ApplyToCallContext(ctx context.Context, container googlecloud.ResourceContainer) context.Context {
	return metadata.AppendToOutgoingContext(ctx, iamTokenKey, i.token)
}

// ApplyToRawHTTPHeader implements googlecloud.CallOptionInjectorOption.
// It sets the IAM token in the HTTP header for REST calls.
func (i *iamTokenCallOptionInjectorOption) ApplyToRawHTTPHeader(header http.Header, container googlecloud.ResourceContainer) {
	header.Set(iamTokenKey, i.token)
}

var _ googlecloud.CallOptionInjectorOption = (*iamTokenCallOptionInjectorOption)(nil)
