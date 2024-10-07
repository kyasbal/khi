package httpclient

import (
	"context"
	"net/http"
)

// An interface to mock *http.Client
type HttpClient[T any] interface {
	// DoWithContext send request to
	DoWithContext(ctx context.Context, request *http.Request) (T, error)
}
