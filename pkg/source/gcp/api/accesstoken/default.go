package accesstoken

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

var MinWaitTimeOnRetriableError = 5
var MaxWaitTimeOnRetriableError = 60
var MaxRetryCount = 3
var RetriableHttpResponseCodes = []int{
	429, 500, 501, 502, 503,
}

var DefaultAccessTokenStore = token.NewBasicTokenStore(
	"accesstoken", token.NewMultiTokenResolver(
		token.NewEnvironmentVariableTokenResolver("GCP_ACCESS_TOKEN"),
		NewMetadataServerAccessTokenResolver(httpclient.NewJsonResponseHttpClient[MDSResponse](httpclient.NewRetryHttpClient(httpclient.NewBasicHttpClient(), MinWaitTimeOnRetriableError, MaxWaitTimeOnRetriableError, MaxRetryCount, RetriableHttpResponseCodes, []int{}, &token.NopTokenRefresher{}, &httpclient.NopTokenApplier{}))),
		&GCloudCommandAccessTokenResolver{},
	),
)
