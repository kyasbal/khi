package quotaproject

import (
	"net/http"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/httpclient"
)

// GCPQuotaProjectHeaderProvider is an implementation of HTTPHeaderProvider for setting quota project.
type GCPQuotaProjectHeaderProvider struct {
	QuotaProject string
}

func NewHeaderProvider(quotaProject string) *GCPQuotaProjectHeaderProvider {
	return &GCPQuotaProjectHeaderProvider{
		QuotaProject: quotaProject,
	}
}

// AddHeader implements httpclient.HTTPHeaderProvider.
func (g *GCPQuotaProjectHeaderProvider) AddHeader(req *http.Request) error {
	if g.QuotaProject == "" {
		return nil
	}
	req.Header.Set("X-Goog-User-Project", g.QuotaProject)
	return nil
}

var _ httpclient.HTTPHeaderProvider = (*GCPQuotaProjectHeaderProvider)(nil)
