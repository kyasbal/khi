package iamtoken

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
)

var DefaultIAMTokenStore = token.NewBasicTokenStore("iamtoken", token.NewMultiTokenResolver(
	token.NewEnvironmentVariableTokenResolver("IAM_TOKEN"),
	&PopupIAMTokenResolver{},
))
