package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

type mockFailClient struct {
	Responses    []*http.Response
	Requests     []*http.Request
	RequestCount int
}

type tokenApplierClientSpy struct {
	CallCount int
}

// ApplyCurrentToken implements TokenApplier.
func (t *tokenApplierClientSpy) ApplyCurrentToken(ctx context.Context, req *http.Request) (*TokenApplyResult, error) {
	t.CallCount++
	return &TokenApplyResult{
		TokenObtainedAt: time.Now(),
	}, nil
}

var _ TokenApplier = (*tokenApplierClientSpy)(nil)

type tokenRefresherClientSpy struct {
	CallCount int
}

// Refresh implements TokenRefresher.
func (t *tokenRefresherClientSpy) Refresh(ctx context.Context) error {
	t.CallCount++
	return nil
}

var _ token.TokenRefresher = (*tokenRefresherClientSpy)(nil)

// DoWithContext implements HttpClient.
func (m *mockFailClient) DoWithContext(ctx context.Context, request *http.Request) (*http.Response, error) {
	m.RequestCount += 1
	m.Requests = append(m.Requests, request)
	return m.Responses[m.RequestCount-1], nil
}

var _ HttpClient[*http.Response] = (*mockFailClient)(nil)

func TestIsRetriable(t *testing.T) {
	type testCase struct {
		RetriableHttpCodes []int
		HttpCode           int
		Expected           bool
	}
	testCases := []testCase{
		{
			RetriableHttpCodes: []int{400, 401, 402},
			HttpCode:           400,
			Expected:           true,
		},
		{
			RetriableHttpCodes: []int{401, 402},
			HttpCode:           400,
			Expected:           false,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with codes:%v", tc.RetriableHttpCodes), func(t *testing.T) {
			client := &RetryHttpClient{RetriableHttpCodes: tc.RetriableHttpCodes}
			actual := client.isRetriable(tc.HttpCode)
			if actual != tc.Expected {
				t.Errorf("got retriability %t, want %t", actual, tc.Expected)
			}
		})
	}
}

func TestRetryBehavior(t *testing.T) {
	type testCase struct {
		Title                       string
		ResponseCodes               []int
		RequestBody                 string
		ExpectedRequestCount        int
		ExpectedError               string
		MinWaitTime                 int
		MaxWaitTime                 int
		MaxRetryCount               int
		ExpectedLastCurrentWaitTime int
		ExpectedTokenApplierCall    int
		ExpectedTokenRefresherCall  int
	}
	testCases := []testCase{
		{
			Title:                       "Simple success",
			ResponseCodes:               []int{200},
			RequestBody:                 "foo",
			ExpectedRequestCount:        1,
			ExpectedError:               "",
			MaxRetryCount:               3,
			MinWaitTime:                 1,
			MaxWaitTime:                 4,
			ExpectedLastCurrentWaitTime: 1,
			ExpectedTokenApplierCall:    1,
			ExpectedTokenRefresherCall:  0,
		},
		{
			Title:                       "Non retriable",
			ResponseCodes:               []int{500},
			RequestBody:                 "foo",
			ExpectedRequestCount:        1,
			ExpectedError:               "unretriable error returned(500):\nBODY:",
			MaxRetryCount:               3,
			MinWaitTime:                 1,
			MaxWaitTime:                 4,
			ExpectedLastCurrentWaitTime: 1,
			ExpectedTokenApplierCall:    1,
			ExpectedTokenRefresherCall:  0,
		},
		{
			Title:                       "Multiple retries",
			ResponseCodes:               []int{400, 400, 200},
			RequestBody:                 "foo",
			ExpectedRequestCount:        3,
			ExpectedError:               "",
			MaxRetryCount:               3,
			MinWaitTime:                 1,
			MaxWaitTime:                 4,
			ExpectedLastCurrentWaitTime: 1,
			ExpectedTokenApplierCall:    3,
			ExpectedTokenRefresherCall:  0,
		},
		{
			Title:                       "Multiple retries and exceed maximum",
			ResponseCodes:               []int{400, 400, 400},
			RequestBody:                 "foo",
			ExpectedRequestCount:        3,
			ExpectedError:               "maximum retry count exceeded 3\nStatus codes:[400 400 400]",
			MaxRetryCount:               3,
			MinWaitTime:                 1,
			MaxWaitTime:                 3,
			ExpectedLastCurrentWaitTime: 3,
			ExpectedTokenApplierCall:    3,
			ExpectedTokenRefresherCall:  0,
		},
		{
			Title:                       "Wait time should be increased as exponential",
			ResponseCodes:               []int{400, 400},
			RequestBody:                 "foo",
			ExpectedRequestCount:        2,
			ExpectedError:               "maximum retry count exceeded 2\nStatus codes:[400 400]",
			MaxRetryCount:               2,
			MinWaitTime:                 1,
			MaxWaitTime:                 10,
			ExpectedLastCurrentWaitTime: 4,
			ExpectedTokenApplierCall:    2,
			ExpectedTokenRefresherCall:  0,
		},
		{
			Title:                       "Refresh token when response code require refreshing token",
			ResponseCodes:               []int{401, 200},
			RequestBody:                 "foo",
			ExpectedRequestCount:        2,
			ExpectedError:               "",
			MaxRetryCount:               2,
			MinWaitTime:                 1,
			MaxWaitTime:                 10,
			ExpectedLastCurrentWaitTime: 1,
			ExpectedTokenApplierCall:    2,
			ExpectedTokenRefresherCall:  1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			responses := []*http.Response{}
			for _, respCode := range tc.ResponseCodes {
				responses = append(responses, &http.Response{
					StatusCode: respCode,
				})
			}
			baseClient := mockFailClient{
				Responses: responses,
				Requests:  make([]*http.Request, 0),
			}
			refresherSpy := tokenRefresherClientSpy{}
			applierSpy := tokenApplierClientSpy{}
			retryClient := NewRetryHttpClient(&baseClient, tc.MinWaitTime, tc.MaxWaitTime, tc.MaxRetryCount, []int{400}, []int{401}, &refresherSpy, &applierSpy)
			retryClient.timeUnit = time.Millisecond
			req, err := http.NewRequest("GET", "https://google.com", bytes.NewBuffer([]byte(tc.RequestBody)))
			if err != nil {
				t.Errorf("got error %v, want nil", err)
			}
			response, err := retryClient.DoWithContext(context.Background(), req)
			if tc.ExpectedError == "" {
				if response == nil {
					t.Error("got nil, want response")
				}
				if err != nil {
					t.Errorf("got error %v, want nil", err)
				}
				if baseClient.RequestCount != tc.ExpectedRequestCount {
					t.Errorf("got retry count %d, want %d", baseClient.RequestCount, tc.ExpectedRequestCount)
				}
			} else {
				if err.Error() != tc.ExpectedError {
					t.Errorf("got error %s, want %s", err.Error(), tc.ExpectedError)
				}
				if baseClient.RequestCount != tc.ExpectedRequestCount {
					t.Errorf("got retry count %d, want %d", baseClient.RequestCount, tc.ExpectedRequestCount)
				}
			}
			for _, req := range baseClient.Requests {
				requestBody, err := io.ReadAll(req.Body)
				if err != nil {
					t.Errorf("got error %v, wnt nil", err)
				}
				requestBodyStr := string(requestBody)
				if requestBodyStr != tc.RequestBody {
					t.Errorf("got requestBody %s, want %s", requestBody, tc.RequestBody)
				}
			}
			if tc.ExpectedLastCurrentWaitTime != retryClient.currentWaitSeconds {
				t.Errorf("got wait time %d, want %d", retryClient.currentWaitSeconds, tc.ExpectedLastCurrentWaitTime)
			}
			if tc.ExpectedTokenApplierCall != applierSpy.CallCount {
				t.Errorf("got token applier call count %d, want %d", applierSpy.CallCount, tc.ExpectedTokenApplierCall)
			}
			if tc.ExpectedTokenRefresherCall != refresherSpy.CallCount {
				t.Errorf("got token refresher call count %d, want %d", refresherSpy.CallCount, tc.ExpectedTokenRefresherCall)
			}
		})
	}
}
