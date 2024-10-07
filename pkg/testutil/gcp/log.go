package gcp_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/api/iamtoken"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api/accesstoken"
)

var accessToken string = ""
var accessTokenOnce sync.Once = sync.Once{}

func tryInitToken() {
	if accessToken == "" {
		fromEnv, exist := os.LookupEnv("ACCESS_TOKEN")
		if exist {
			accessToken = fromEnv
		} else {
			slog.Info("Failed to get access token from environment variable. Getting from gcloud command...")
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd := exec.Command("gcloud", "auth", "print-access-token")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				panic(fmt.Sprintf("Failed to get access token via gcloud command. \nstderr:\n%s\n\nstdout:\n%s\n\nerr:%s", stderr.String(), stdout.String(), err.Error()))
			}
			accessToken = strings.ReplaceAll(stdout.String(), "\n", "")
		}
	}
}

func IsValidLogQuery(query string) error {
	accessTokenOnce.Do(tryInitToken)
	gcpApi, err := api.NewGCPClient(accesstoken.DefaultAccessTokenStore, iamtoken.DefaultIAMTokenStore, "")
	if err != nil {
		return err
	}
	query = fmt.Sprintf(`%s
timestamp >= "2024-01-01T00:00:00Z"
timestamp <= "2024-01-01T00:00:01Z"`, query)

	return gcpApi.ListLogEntries(context.Background(), "kubernetes-history-inspector", query, make(chan any))
}
