package grouper

import (
	common_grouper "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/grouper"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
)

type LogGrouper = common_grouper.Grouper[*log.LogEntity, string]
