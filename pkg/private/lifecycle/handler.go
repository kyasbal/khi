package lifecycle

import (
	"os"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/lifecycle"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/analytics"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/analytics/types"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/api/iamtoken"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/private/parameters"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api"
)

// NewAnalyticsLifecycleHandler returns a new LifecycleEventHandler to report analytics events on KHI lifecycle event.
func NewAnalyticsLifecycleHandler() *lifecycle.LifecycleEventHandler {
	return &lifecycle.LifecycleEventHandler{
		OnInit: func() {
			reporter := analytics.NewAnalyticsReporter()
			reporter.ReportEvent(types.AnalyticsEventKHIStart, map[string]any{})
		},
		OnTerminate: func(s os.Signal) {
			reporter := analytics.NewAnalyticsReporter()
			reporter.ReportEvent(types.AnalyticsEventKHITerminate, map[string]any{
				"signal": s.String(),
			})
		},
		OnInspectionStart: func(runId, inspectionType string) {
			reporter := analytics.NewAnalyticsReporter()
			reporter.ReportEvent(types.AnalyticsEventInspectionStart, map[string]any{
				"rid":            runId,
				"inspectionType": inspectionType,
			})
		},
		OnInspectionEnd: func(runId, inspectionType, status string, size int) {
			reporter := analytics.NewAnalyticsReporter()
			reporter.ReportEvent(types.AnalyticsEventInspectionEnd, map[string]any{
				"rid":            runId,
				"inspectionType": inspectionType,
				"status":         status,
				"resultSize":     size,
			})
		},
	}
}

func NewIAMTokenSetupLifecycleHandler() *lifecycle.LifecycleEventHandler {
	return &lifecycle.LifecycleEventHandler{
		OnInit: func() {
			if *parameters.Private.InspectionMode {
				api.DefaultGCPClientFactory.RegisterHeaderProvider(iamtoken.NewHeaderProvider(iamtoken.DefaultIAMTokenStore))
				api.DefaultGCPClientFactory.RegisterRefreshableTokenStore(iamtoken.DefaultIAMTokenStore)
			}
		},
	}

}
