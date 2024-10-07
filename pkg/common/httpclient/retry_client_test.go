package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

type mockFailClient struct {
	Responses    []*http.Response
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
func (t *tokenRefresherClientSpy) Refresh(ctx context.Context) {
	t.CallCount++
}

var _ TokenRefresher = (*tokenRefresherClientSpy)(nil)

// DoWithContext implements HttpClient.
func (m *mockFailClient) DoWithContext(ctx context.Context, request *http.Request) (*http.Response, error) {
	m.RequestCount += 1
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
				t.Errorf("unmatched result. Expected:%t, Actual:%t", tc.Expected, actual)
			}
		})
	}
}

func TestRetryBehavior(t *testing.T) {
	type testCase struct {
		Title                       string
		ResponseCodes               []int
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
			}
			refresherSpy := tokenRefresherClientSpy{}
			applierSpy := tokenApplierClientSpy{}
			retryClient := NewRetryHttpClient(&baseClient, tc.MinWaitTime, tc.MaxWaitTime, tc.MaxRetryCount, []int{400}, []int{401}, &refresherSpy, &applierSpy)
			retryClient.timeUnit = time.Millisecond
			response, err := retryClient.DoWithContext(context.Background(), &http.Request{})
			if tc.ExpectedError == "" {
				if response == nil {
					t.Errorf("response was unexpected nil")
				}
				if err != nil {
					t.Errorf("unexpected error %s", err.Error())
				}
				if baseClient.RequestCount != tc.ExpectedRequestCount {
					t.Errorf("unexpected retry count, expected %d, but %d", tc.ExpectedRequestCount, baseClient.RequestCount)
				}
			} else {
				if err.Error() != tc.ExpectedError {
					t.Errorf("unexpected error %s, expected %s", err.Error(), tc.ExpectedError)
				}
				if baseClient.RequestCount != tc.ExpectedRequestCount {
					t.Errorf("unexpected retry count, expected %d, but %d", tc.ExpectedRequestCount, baseClient.RequestCount)
				}
			}
			if tc.ExpectedLastCurrentWaitTime != retryClient.currentWaitSeconds {
				t.Errorf("unexpected wait time %d, expected %d", retryClient.currentWaitSeconds, tc.ExpectedLastCurrentWaitTime)
			}
			if tc.ExpectedTokenApplierCall != applierSpy.CallCount {
				t.Errorf("unexpected token applier call count %d, expected %d", applierSpy.CallCount, tc.ExpectedTokenApplierCall)
			}
			if tc.ExpectedTokenRefresherCall != refresherSpy.CallCount {
				t.Errorf("unexpected token refresher call count %d, expected %d", refresherSpy.CallCount, tc.ExpectedTokenRefresherCall)
			}
		})
	}
}
