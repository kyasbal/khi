// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lifecycle

import (
	"os"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/errorreport"
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

// NewErrorReportLifecycleHandler returns a new LifecycleEventHandler to register GA metadata labels to the error reporter.
func NewErrorReportLifecycleHandler() *lifecycle.LifecycleEventHandler {
	return &lifecycle.LifecycleEventHandler{
		OnInit: func() {
			if parameters.Private.GALabels != nil {
				metadata := parameters.Private.GetMapOfGALabels()
				for key, value := range metadata {
					errorreport.DefaultErrorReporter.SetMetadataEntry(key, value)
				}
			}
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
