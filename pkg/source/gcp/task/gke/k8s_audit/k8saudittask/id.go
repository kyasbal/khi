package k8saudittask

import (
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
)

const K8sAuditQueryTaskId = gcp_task.GCPPrefix + "query/k8s_audit"
const K8sAuditParseTaskId = gcp_task.GCPPrefix + "/feature/audit-parser-v2"
const k8sAuditTaskIDPrefix = gcp_task.GCPPrefix + "feature/k8s_audit/"

const TimelineGroupingTaskId = k8sAuditTaskIDPrefix + "timelne-grouping"
const ManifestGenerateTaskId = k8sAuditTaskIDPrefix + "manifest-generate"
const LogConvertTaskId = k8sAuditTaskIDPrefix + "log-convert"
const CommonLogParseTaskId = k8sAuditTaskIDPrefix + "common-fields-parse"
