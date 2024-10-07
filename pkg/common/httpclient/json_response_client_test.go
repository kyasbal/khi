package httpclient

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type mockHttpClient struct {
	Response *http.Response
	Error    error
}

// DoWithContext implements HttpClient.
func (m *mockHttpClient) DoWithContext(ctx context.Context, request *http.Request) (*http.Response, error) {
	if m.Error == nil {
		return m.Response, nil
	} else {
		return nil, m.Error
	}
}

var _ HttpClient[*http.Response] = (*mockHttpClient)(nil)

func TestDoWithContext(t *testing.T) {
	type testJsonType struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}
	jsonClient := NewJsonResponseHttpClient[testJsonType](&mockHttpClient{
		Response: &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
  "foo":"foo-val",
  "bar":"bar-val"
}`)),
		},
	})
	result, _, err := jsonClient.DoWithContext(context.Background(), &http.Request{})
	if err != nil {
		t.Errorf("unexpected err:%s", err.Error())
	}
	if diff := cmp.Diff(&testJsonType{
		Foo: "foo-val",
		Bar: "bar-val",
	}, result); diff != "" {
		t.Errorf("response is not matching with the expected value\n%s", diff)
	}
}
