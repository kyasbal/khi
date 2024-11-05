package iamtoken

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/token"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parameters"
)

var DefaultIAMTokenStore = token.NewBasicTokenStore("iamtoken", token.NewMultiTokenResolver(
	token.NewOnceTokenResolver(func() string {
		if parameters.Private.IAMToken == nil {
			return ""
		}
		return *parameters.Private.IAMToken
	}),
	&PopupIAMTokenResolver{},
))
