package accesstoken

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
)

var MinWaitTimeOnRetriableError = 5
var MaxWaitTimeOnRetriableError = 60
var MaxRetryCount = 3
var RetriableHttpResponseCodes = []int{
	429, 500, 501, 502, 503,
}

var DefaultOAuthTokenResolver = NewOAuthTokenResolver()

var DefaultAccessTokenStore = token.NewBasicTokenStore(
	"accesstoken", token.NewMultiTokenResolver(
		DefaultOAuthTokenResolver,
		token.NewOnceTokenResolver(func() string {
			if parameters.Auth.AccessToken == nil {
				return ""
			}
			return *parameters.Auth.AccessToken
		}),
		NewMetadataServerAccessTokenResolver(httpclient.NewJsonResponseHttpClient[MDSResponse](httpclient.NewRetryHttpClient(httpclient.NewBasicHttpClient(), MinWaitTimeOnRetriableError, MaxWaitTimeOnRetriableError, MaxRetryCount, RetriableHttpResponseCodes, []int{}, &token.NopTokenRefresher{}, &httpclient.NopTokenApplier{}))),
		&GCloudCommandAccessTokenResolver{},
	),
)
