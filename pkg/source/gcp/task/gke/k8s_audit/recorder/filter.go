package recorder

import (
	"context"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task/gke/k8s_audit/types"
)

func AnyLogGroupFilter() LogGroupFilterFunc {
	return func(ctx context.Context, resourcePath string) bool {
		return true
	}
}

func ResourceKindLogGroupFilter(kindInSingular string) LogGroupFilterFunc {
	return func(ctx context.Context, resourcePath string) bool {
		pathSegments := strings.Split(resourcePath, "#")
		if len(pathSegments) == 4 {
			return kindInSingular == pathSegments[1]
		} else {
			return false
		}
	}
}

func SubresourceLogGroupFilter(subresource string) LogGroupFilterFunc {
	return func(ctx context.Context, resourcePath string) bool {
		pathSegments := strings.Split(resourcePath, "#")
		if len(pathSegments) == 5 {
			return subresource == pathSegments[4]
		} else {
			return false
		}
	}
}

func AnyLogFilter() LogFilterFunc {
	return func(ctx context.Context, l *types.ResourceSpecificParserInput) bool {
		return true
	}
}

// OnlySucceedLogs returns a LogFilterFunc that only matches audit logs with non zero response code.
func OnlySucceedLogs() LogFilterFunc {
	return func(ctx context.Context, l *types.ResourceSpecificParserInput) bool {
		return l.Code == 0
	}
}

func OnlyWithResourceBody() LogFilterFunc {
	return func(ctx context.Context, l *types.ResourceSpecificParserInput) bool {
		return l.ResourceBodyReader != nil
	}
}

func OnlySpecificVerb(verb enum.RevisionVerb) LogFilterFunc {
	return func(ctx context.Context, l *types.ResourceSpecificParserInput) bool {
		return l.Operation.Verb == verb
	}
}

func AndLogFilter(filters ...LogFilterFunc) LogFilterFunc {
	return func(ctx context.Context, l *types.ResourceSpecificParserInput) bool {
		for _, filter := range filters {
			if !filter(ctx, l) {
				return false
			}
		}
		return true
	}
}
