package gcp_test

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/api/iamtoken"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api/accesstoken"
)

func IsValidLogQuery(query string) error {
	gcpApi, err := api.NewGCPClient(accesstoken.DefaultAccessTokenStore, iamtoken.DefaultIAMTokenStore, "")
	if err != nil {
		return err
	}
	query = fmt.Sprintf(`%s
timestamp >= "2024-01-01T00:00:00Z"
timestamp <= "2024-01-01T00:00:01Z"`, query)

	return gcpApi.ListLogEntries(context.Background(), "kubernetes-history-inspector", query, make(chan any))
}
