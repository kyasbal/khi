package quotaproject

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGCPQuotaProjectHeaderProvider(t *testing.T) {
	testCases := []struct {
		name         string
		quotaProject string
		wantHeader   string
		wantErr      bool
	}{
		{
			name:         "success",
			quotaProject: "test-project",
			wantHeader:   "test-project",
			wantErr:      false,
		},
		{
			name:         "empty quota project",
			quotaProject: "",
			wantHeader:   "",
			wantErr:      false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p := NewHeaderProvider(tt.quotaProject)
			req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			err = p.AddHeader(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GCPQuotaProjectHeaderProvider.AddHeader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				gotHeader := req.Header.Get("X-Goog-User-Project")
				if diff := cmp.Diff(tt.wantHeader, gotHeader); diff != "" {
					t.Errorf("GCPQuotaProjectHeaderProvider.AddHeader() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
