package httpclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("should return error when server is down", func(t *testing.T) {
		client := NewBasicHttpClient()
		req, _ := http.NewRequest("GET", "http://localhost:12345", nil)

		_, err := client.DoWithContext(context.Background(), req)

		if err == nil {
			t.Error("Expected error, but got nil")
		}
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

		if err == nil {
			t.Error("Expected error, but got nil")
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected error to be context.Canceled, but got %v", err)
		}
	})
}
