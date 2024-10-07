package api

import (
	"context"
	"net/http"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
)

type RequestGenerator = func(hasToken bool, nextPageToken string) (*http.Request, error)

// PageClient is utility to obtain all the resource from API returning page token.
type PageClient[T any] struct {
	client httpclient.HttpClient[*http.Response]
}

func NewPageClient[T any](client httpclient.HttpClient[*http.Response]) *PageClient[T] {
	return &PageClient[T]{
		client: client,
	}
}

func (p *PageClient[T]) GetAll(ctx context.Context, requestGenerator RequestGenerator, nextPageTokenMapper func(response *T) string) ([]*T, error) {
	result := make([]*T, 0)
	for nextPageToken := "-"; nextPageToken != ""; {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			request, err := requestGenerator(nextPageToken != "-", nextPageToken)
			if err != nil {
				return nil, err
			}
			client := httpclient.NewJsonResponseHttpClient[T](p.client)
			typedResponse, _, err := client.DoWithContext(ctx, request)
			if err != nil {
				return nil, err
			}
			result = append(result, typedResponse)
			nextPageToken = nextPageTokenMapper(typedResponse)
		}
	}
	return result, nil
}
