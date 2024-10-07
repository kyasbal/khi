package types

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/rtype"
)

// ResourceSpecificParserInput is a type passed to ResourceSpecificParser from the prestep parser.
type ResourceSpecificParserInput struct {
	// Current Log
	Log *log.LogEntity
	// ResourceName field of the log
	ResourceName string
	// MethodName field of the log
	MethodName string
	// PrincipalEmail field of the log
	PrincipalEmail string
	// Kubernetes operation read from resource name and method name
	Operation *model.KubernetesObjectOperation
	// The request field of this log. This can be nil depending on the audit policy.
	Request     *structure.Reader
	RequestType rtype.Type
	// The response field of this log. This can be nil depending on the audit policy.
	Response     *structure.Reader
	ResponseType rtype.Type
	// Current resource body canged by this request.
	ResourceBodyReader *structure.Reader
	ResourceBodyYaml   string
	// The response code.
	Code int

	GeneratedFromDeleteCollectionOperation bool
}

type TimelineGrouperResult struct {
	TimelineResourcePath string
	PreParsedLogs        []*ResourceSpecificParserInput
}
