package iamtoken

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/popup"
)

type IAMTokenResolverPopupForm struct {
}

// GetMetadata implements popup.PopupForm.
func (i *IAMTokenResolverPopupForm) GetMetadata() popup.PopupFormMetadata {
	return popup.PopupFormMetadata{
		Title:       "Google Admin Token",
		Description: "Refresh Google Admin token by pasting the KHI command again on this form",
		Type:        "text",
		Placeholder: "gcloud auth configure-docker --quiet\n" +
			"docker run --rm --pull always -p 8080:8080 -it -e KHI_FIXED_PROJECT_ID=\"XXXXXXXXX\" -e GCP_DEFAULT_PROJECT=\"XXXXXXXXX\" -e GCP_ACCESS_TOKEN=`gcloud auth print-access-token` -e GCP_IDENTITY_TOKEN=`gcloud auth print-identity-token` -e KHI_GA_LABELS=\"justification=vector/123456,user=XXXXXXXXX\" -e IAM_TOKEN=\"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\\n" +
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\\n" +
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\\n" +
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\\n" +
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\\n" +
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\\n" +
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\\n" +
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\\\n" +
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\" gcr.io/kubernetes-history-inspector/standalone:latest",
	}
}

// Validate implements popup.PopupForm.
func (i *IAMTokenResolverPopupForm) Validate(req *popup.PopupAnswerResponse) string {
	_, err := extractTokenFromGACommand(req.Value)
	if err != nil {
		return err.Error()
	}
	return ""
}

var _ popup.PopupForm = (*IAMTokenResolverPopupForm)(nil)

type PopupIAMTokenResolver struct {
}

// Resolve implements token.TokenResolver.
func (p *PopupIAMTokenResolver) Resolve(ctx context.Context) (*token.Token, error) {
	_, found := os.LookupEnv("USE_IAM_TOKEN")
	if !found {
		return nil, nil
	}
	slog.InfoContext(ctx, "Requesting a new token with showing a popup.")
	gaCommand, err := popup.Instance.ShowPopup(&IAMTokenResolverPopupForm{})
	if err != nil {
		return nil, err
	}
	slog.InfoContext(ctx, "New IAM token received from the popup.")
	rawToken, err := extractTokenFromGACommand(gaCommand)
	if err != nil {
		return nil, err
	}
	return token.New(rawToken), nil
}

var _ token.TokenResolver = (*PopupIAMTokenResolver)(nil)

func extractTokenFromGACommand(command string) (string, error) {
	splitted := strings.Split(command, " ")
	for _, segment := range splitted {
		if strings.HasPrefix(segment, "IAM_TOKEN=") {
			return strings.ReplaceAll(strings.TrimSuffix(strings.TrimPrefix(segment, "IAM_TOKEN=\""), "\""), "\\\n", ""), nil
		}
	}
	return "", fmt.Errorf("failed to parse given command. IAM_TOKEN section not found from the input. Did you input the KHI command shown on Google Admin correctly?")
}
