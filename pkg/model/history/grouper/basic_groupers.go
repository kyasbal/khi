package grouper

import (
	common_grouper "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/grouper"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
)

var AllIndependentLogGrouper LogGrouper = common_grouper.NewBasicGrouper[*log.LogEntity, string](func(log *log.LogEntity) string {
	return log.ID()
})

var AllDependentLogGrouper LogGrouper = common_grouper.NewBasicGrouper[*log.LogEntity, string](func(log *log.LogEntity) string {
	return ""
})

func NewSingleStringFieldKeyLogGrouper(keyPath string) LogGrouper {
	return common_grouper.NewBasicGrouper[*log.LogEntity, string](func(log *log.LogEntity) string {
		return log.GetStringOrDefault(keyPath, "")
	})
}
