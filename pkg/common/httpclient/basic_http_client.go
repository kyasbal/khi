package httpclient

import (
	"context"
	"net/http"
)

type BasicHttpClient struct {
}

// BasicHttpClient implements HttpClient interface
var _ HttpClient[*http.Response] = (*BasicHttpClient)(nil)

// DoWithContext implements HttpClient.
func (b *BasicHttpClient) DoWithContext(ctx context.Context, request *http.Request) (*http.Response, error) {
	req := request.WithContext(ctx)
	client := new(http.Client)
	return client.Do(req)
}

func NewBasicHttpClient() HttpClient[*http.Response] {
	return &BasicHttpClient{}
}
