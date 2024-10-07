package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JSONReponseHttpClient[T any] struct {
	client HttpClient[*http.Response]
}

func NewJsonResponseHttpClient[T any](client HttpClient[*http.Response]) *JSONReponseHttpClient[T] {
	return &JSONReponseHttpClient[T]{
		client: client,
	}
}

func (p *JSONReponseHttpClient[T]) DoWithContext(ctx context.Context, request *http.Request) (*T, *http.Response, error) {
	response, err := p.client.DoWithContext(ctx, request)
	if err != nil {
		return nil, response, err
	}
	if response.StatusCode >= 400 {
		return nil, response, fmt.Errorf("%d:%s", response.StatusCode, response.Status)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, response, err
	}
	result, err := p.parse(string(responseData))
	if err != nil {
		return nil, response, err
	}
	return result, response, err
}

func (p *JSONReponseHttpClient[T]) parse(body string) (*T, error) {
	if body == "" {
		return nil, fmt.Errorf("response is empty")
	}
	var typedResponse T
	err := json.Unmarshal([]byte(body), &typedResponse)
	if err != nil {
		return nil, err
	}
	return &typedResponse, nil
}
