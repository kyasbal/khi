package api

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api/accesstoken"
)

type GCPClientFactory struct {
	HeaderProviders []httpclient.HTTPHeaderProvider
	TokenStores     []token.TokenStore
}

func NewGCPClientFactory() *GCPClientFactory {
	return &GCPClientFactory{
		TokenStores:     []token.TokenStore{},
		HeaderProviders: []httpclient.HTTPHeaderProvider{},
	}
}

// NewClient instanciate a new GCPClient from current factory config.
func (f *GCPClientFactory) NewClient() (GCPClient, error) {
	return NewGCPClient(token.NewMultiTokenStoreRefresher(f.TokenStores...), f.HeaderProviders)
}

// RegisterHeaderProvider adds a new HeaderProvider on factory config.
func (f *GCPClientFactory) RegisterHeaderProvider(provider httpclient.HTTPHeaderProvider) {
	f.HeaderProviders = append(f.HeaderProviders, provider)
}

// RegisterRefreshableTokenStore adds a refreshable token store. The token will be refreshed when permission related error happens.
func (f *GCPClientFactory) RegisterRefreshableTokenStore(store token.TokenStore) {
	f.TokenStores = append(f.TokenStores, store)
}

var DefaultGCPClientFactory *GCPClientFactory = NewGCPClientFactory()

// set the default header providers on the default factory.
func init() {
	DefaultGCPClientFactory.RegisterHeaderProvider(accesstoken.NewHeaderProvider(accesstoken.DefaultAccessTokenStore))
	DefaultGCPClientFactory.RegisterRefreshableTokenStore(accesstoken.DefaultAccessTokenStore)
}
