package k8s_event

import (
	"context"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/grouper"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history/resourcepath"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parser"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var GKEK8sEventLogParseJob = parser.NewParserTaskFromParser(gcp_task.GCPPrefix+"feature/event-parser", &k8sEventParser{}, true)

type k8sEventParser struct {
}

// Description implements parser.Parser.
func (*k8sEventParser) Description() string {
	return `Visualize Kubernetes event logs on GKE.
This parser shows events associated to K8s resources`
}

// GetParserName implements parser.Parser.
func (*k8sEventParser) GetParserName() string {
	return `Kubernetes Event Logs`
}

func (*k8sEventParser) Dependencies() []string {
	return []string{}
}

func (*k8sEventParser) LogTask() string {
	return GKEK8sEventLogQueryTaskId
}

func (*k8sEventParser) Grouper() grouper.LogGrouper {
	return grouper.AllDependentLogGrouper
}

// Parse implements parser.Parser.
func (*k8sEventParser) Parse(ctx context.Context, l *log.LogEntity, cs *history.ChangeSet, builder *history.Builder, v *task.VariableSet) error {
	if kind, err := l.GetString("jsonPayload.kind"); err != nil {
		// Event exporter ingests cluster scoped logs without jsonPayload
		if textPayload, err := l.GetString("textPayload"); err == nil {
			clusterName := l.GetStringOrDefault("resource.labels.cluster_name", "Unknown")
			cs.RecordEvent(resourcepath.Cluster(clusterName))
			cs.RecordLogSummary(textPayload)
			return nil
		}
		return err
	} else {
		if kind != "Event" {
			return fmt.Errorf("skipping kind:%s", kind)
		}
	}
	apiVersion := l.GetStringOrDefault("jsonPayload.involvedObject.apiVersion", "v1")

	kind := l.GetStringOrDefault("jsonPayload.involvedObject.kind", "Unknown")

	name := l.GetStringOrDefault("jsonPayload.involvedObject.name", "Unknown")

	namespace := l.GetStringOrDefault("jsonPayload.involvedObject.namespace", "cluster-scope")
	if !strings.Contains(apiVersion, "/") {
		apiVersion = "core/" + apiVersion
	}

	cs.RecordEvent(resourcepath.NameLayerGeneralItem(apiVersion, strings.ToLower(kind), namespace, name))
	cs.RecordLogSummary(fmt.Sprintf("【%s】%s", l.GetStringOrDefault("jsonPayload.reason", "Unknown"), l.GetStringOrDefault("jsonPayload.message", "")))
	return nil
}

var _ parser.Parser = (*k8sEventParser)(nil)
