package httpclient_test

import (
	"context"
	"net/http"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
)

type HttpClientSpyResponse[T any] struct {
	Response T
	Error    error
}

type HttpClientSpy[T any] struct {
	Results  []*HttpClientSpyResponse[T]
	Requests []*http.Request
}

func NewHttpClientSpyResponse[T any](response T, err error) *HttpClientSpyResponse[T] {
	return &HttpClientSpyResponse[T]{
		Response: response,
		Error:    err,
	}
}

func NewHttpClientSpy[T any](responses ...*HttpClientSpyResponse[T]) *HttpClientSpy[T] {
	return &HttpClientSpy[T]{
		Results:  responses,
		Requests: make([]*http.Request, 0),
	}
}

// DoWithContext implements httpclient.HttpClient.
func (h *HttpClientSpy[T]) DoWithContext(ctx context.Context, request *http.Request) (T, error) {
	h.Requests = append(h.Requests, request)
	callIndex := len(h.Requests) - 1
	if callIndex >= len(h.Results) {
		return *new(T), nil
	}
	return h.Results[callIndex].Response, h.Results[callIndex].Error
}

var _ httpclient.HttpClient[any] = (*HttpClientSpy[any])(nil)
