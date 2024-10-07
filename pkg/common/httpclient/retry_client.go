package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type RetryHttpClient struct {
	Client                             HttpClient[*http.Response]
	MinWaitSeconds                     int
	MaxWaitSeconds                     int
	MaxRetryCount                      int
	RetriableHttpCodes                 []int
	RetriableWithRefreshTokenHttpCodes []int
	currentWaitSeconds                 int
	timeUnit                           time.Duration // For testing purpose to make test faster
	tokenRefresher                     TokenRefresher
	tokenApplier                       TokenApplier
}

func NewRetryHttpClient(baseClient HttpClient[*http.Response], minWaitSeconds int, maxWaitSeconds int, maxRetryCount int, retriableHttpCodes []int, retriableWithRefreshTokenHttpCodes []int, tokenRefresher TokenRefresher, tokenApplier TokenApplier) *RetryHttpClient {
	return &RetryHttpClient{
		Client:                             baseClient,
		MinWaitSeconds:                     minWaitSeconds,
		MaxWaitSeconds:                     maxWaitSeconds,
		MaxRetryCount:                      maxRetryCount,
		RetriableHttpCodes:                 retriableHttpCodes,
		RetriableWithRefreshTokenHttpCodes: retriableWithRefreshTokenHttpCodes,
		currentWaitSeconds:                 minWaitSeconds,
		timeUnit:                           time.Second,
		tokenRefresher:                     tokenRefresher,
		tokenApplier:                       tokenApplier,
	}
}

// DoWithContext implements HttpClient.
func (r *RetryHttpClient) DoWithContext(ctx context.Context, request *http.Request) (*http.Response, error) {
	// Clone request body into array to create another reader of Body on retry.
	var clonedRequest []byte
	if request.Body != nil {
		var err error
		clonedRequest, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		request.Body = io.NopCloser(bytes.NewBuffer(clonedRequest))
	}
	statusCodes := []int{}
	for i := 0; i < r.MaxRetryCount; i++ {
		_, err := r.tokenApplier.ApplyCurrentToken(ctx, request)
		if err != nil {
			return nil, err
		}
		response, err := r.Client.DoWithContext(ctx, request)
		if err != nil {
			return nil, err
		}
		if response.StatusCode < 400 {
			r.currentWaitSeconds = r.MinWaitSeconds
			// Treat this response is ok not to retry
			return response, nil
		}
		if !r.isRetriable(response.StatusCode) {
			body := []byte{}
			if response.Body != nil {
				body, _ = io.ReadAll(response.Body)
			}
			return response, fmt.Errorf("unretriable error returned(%d):%s\nBODY:%s", response.StatusCode, response.Status, string(body))
		} else {
			statusCodes = append(statusCodes, response.StatusCode)
			if r.isRetriableWithRefreshingToken(response.StatusCode) {
				slog.DebugContext(ctx, fmt.Sprintf("Previous request to %s got %d response. Attempting retrying with refreshing the token.", request.RequestURI, response.StatusCode))
				r.tokenRefresher.Refresh(ctx)
				r.currentWaitSeconds = r.MinWaitSeconds
			} else {
				r.currentWaitSeconds *= 2
				if r.currentWaitSeconds > r.MaxWaitSeconds {
					r.currentWaitSeconds = r.MaxWaitSeconds
				}
				slog.DebugContext(ctx, fmt.Sprintf("Previous request to %s got %d response. Next retry after %d seconds", request.RequestURI, response.StatusCode, r.currentWaitSeconds))
				time.Sleep(r.timeUnit * time.Duration(r.currentWaitSeconds))
				if request.Body != nil {
					request.Body = io.NopCloser(bytes.NewBuffer(clonedRequest))
				}
			}
		}
	}
	return nil, fmt.Errorf("maximum retry count exceeded %d\nStatus codes:%v", r.MaxRetryCount, statusCodes)
}

func (r *RetryHttpClient) isRetriable(code int) bool {
	for _, retryCode := range r.RetriableHttpCodes {
		if code == retryCode {
			return true
		}
	}
	return r.isRetriableWithRefreshingToken(code)
}

func (r *RetryHttpClient) isRetriableWithRefreshingToken(code int) bool {
	for _, retryCode := range r.RetriableWithRefreshTokenHttpCodes {
		if code == retryCode {
			return true
		}
	}
	return false
}

var _ (HttpClient[*http.Response]) = (*RetryHttpClient)(nil)
