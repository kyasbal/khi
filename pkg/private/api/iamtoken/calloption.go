package iamtoken

import (
	"context"
	"net/http"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typeddict"
	"google.golang.org/grpc/metadata"
)

const iamTokenKey = "x-goog-iam-authorization-token"

// IAMTokenCallOptionInjectorOption is a CallOptionInjectorOption that injects an IAM token into API calls.
type IAMTokenCallOptionInjectorOption struct {
	defaultToken            string
	containerSpecificTokens *typeddict.TypedDict[string]
}

func New(defaultToken string) *IAMTokenCallOptionInjectorOption {
	return &IAMTokenCallOptionInjectorOption{
		defaultToken:            defaultToken,
		containerSpecificTokens: typeddict.NewTypedDict[string](),
	}
}

// ApplyToCallContext implements googlecloud.CallOptionInjectorOption.
// It appends the IAM token to the outgoing context metadata for gRPC calls.
func (o *IAMTokenCallOptionInjectorOption) ApplyToCallContext(ctx context.Context, container googlecloud.ResourceContainer) context.Context {
	token := o.getTokenFor(container)
	return metadata.AppendToOutgoingContext(ctx, iamTokenKey, token)
}

// ApplyToRawHTTPHeader implements googlecloud.CallOptionInjectorOption.
// It sets the IAM token in the HTTP header for REST calls.
func (o *IAMTokenCallOptionInjectorOption) ApplyToRawHTTPHeader(header http.Header, container googlecloud.ResourceContainer) {
	token := o.getTokenFor(container)
	header.Set(iamTokenKey, token)
}

// SetTokenFor sets the IAM token for the given container.
func (o *IAMTokenCallOptionInjectorOption) SetTokenFor(container googlecloud.ResourceContainer, token string) {
	typeddict.Set(o.containerSpecificTokens, container.Identifier(), token)
}

func (o *IAMTokenCallOptionInjectorOption) getTokenFor(container googlecloud.ResourceContainer) string {
	if container == nil {
		return o.defaultToken
	}
	ci := container.Identifier()
	token, found := typeddict.Get(o.containerSpecificTokens, ci)
	if !found {
		return o.defaultToken
	}
	return token
}

var _ googlecloud.CallOptionInjectorOption = (*IAMTokenCallOptionInjectorOption)(nil)
