package httpclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasicHttpClient_DoWithContext(t *testing.T) {
	t.Run("should return response when server returns 200", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}))
		defer ts.Close()

		client := NewBasicHttpClient()
		req, _ := http.NewRequest("GET", ts.URL, nil)

		resp, err := client.DoWithContext(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return error when server is down", func(t *testing.T) {
		client := NewBasicHttpClient()
		req, _ := http.NewRequest("GET", "http://localhost:12345", nil)

		_, err := client.DoWithContext(context.Background(), req)

		assert.Error(t, err)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		client := NewBasicHttpClient()
		req, _ := http.NewRequest("GET", ts.URL, nil)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := client.DoWithContext(ctx, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled))
	})
}
