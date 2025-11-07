// Copyright 2025 Google LLC
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

// Package googlecloudlogk8snode_contract defines the contract for the googlecloudlogk8snode task.
package googlecloudlogk8snode_contract

import (
	inspectiontaskbase "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/taskbase"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/model/log"
)

const (
	// TaskIDPrefix is the prefix for all task IDs in this package.
	TaskIDPrefix = "cloud.google.com/log/k8s-node/"
)

// ListLogEntriesTaskID is the task id for the task that queries k8s node logs from Cloud Logging.
var ListLogEntriesTaskID = taskid.NewDefaultImplementationID[[]*log.Log](TaskIDPrefix + "query")

// LogSerializerTaskID is the task ID to finalize the logs to be included in the final output.
var LogSerializerTaskID = taskid.NewDefaultImplementationID[[]*log.Log](TaskIDPrefix + "log-serializer")

var CommonFieldsetReaderTaskID = taskid.NewDefaultImplementationID[[]*log.Log](TaskIDPrefix + "common-fieldset-reader")

// ContainerdLogFilterTaskID is the task ID for filtering containerd logs.
var ContainerdLogFilterTaskID = taskid.NewDefaultImplementationID[[]*log.Log](TaskIDPrefix + "containerd-log-filter")

var ContainerdLogGroupTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](TaskIDPrefix + "containerd-log-group")

var ContainerdIDDiscoveryTaskID = taskid.NewDefaultImplementationID[*ContainerdRelationshipRegistry](TaskIDPrefix + "containerd-id-discovery")

var ContainerdLogHistoryModifierTaskID = taskid.NewDefaultImplementationID[struct{}](TaskIDPrefix + "containerd-log-history-modifier")

var KubeletLogFilterTaskID = taskid.NewDefaultImplementationID[[]*log.Log](TaskIDPrefix + "kubelet-log-filter")

var KubeletLogGroupTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](TaskIDPrefix + "kubelet-log-group")

var KubeletLogHistoryModifierTaskID = taskid.NewDefaultImplementationID[struct{}](TaskIDPrefix + "kubelet-log-history-modifier")

var KubeProxyLogFilterTaskID = taskid.NewDefaultImplementationID[[]*log.Log](TaskIDPrefix + "kube-proxy-log-filter")

var KubeProxyLogGroupTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](TaskIDPrefix + "kube-proxy-log-group")

var KubeProxyLogHistoryModifierTaskID = taskid.NewDefaultImplementationID[struct{}](TaskIDPrefix + "kube-proxy-log-history-modifier")

// OtherLogFilterTaskID is the task ID for filtering other logs.
var OtherLogFilterTaskID = taskid.NewDefaultImplementationID[[]*log.Log](TaskIDPrefix + "other-log-filter")

var OtherLogGroupTaskID = taskid.NewDefaultImplementationID[inspectiontaskbase.LogGroupMap](TaskIDPrefix + "other-log-group")

var OtherLogHistoryModifierTaskID = taskid.NewDefaultImplementationID[struct{}](TaskIDPrefix + "other-log-history-modifier")

var TailTaskID = taskid.NewDefaultImplementationID[struct{}](TaskIDPrefix + "tail")
