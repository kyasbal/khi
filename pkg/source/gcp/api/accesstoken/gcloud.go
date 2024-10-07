package accesstoken

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

type GCloudCommandAccessTokenResolver struct {
}

// Resolve implements token.TokenResolver.
func (g *GCloudCommandAccessTokenResolver) Resolve(ctx context.Context, expiredTokens map[string]interface{}) (string, error) {
	slog.InfoContext(ctx, `Environment variable "GCP_ACCESS_TOKEN" was not found. Trying to get access token with gcloud command...`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get access token via gcloud command. \nstderr:\n%s\n\nstdout:\n%s\n\nerr:%s", stderr.String(), stdout.String(), err.Error())
	}
	return strings.ReplaceAll(stdout.String(), "\n", ""), nil
}

var _ token.TokenResolver = (*GCloudCommandAccessTokenResolver)(nil)
