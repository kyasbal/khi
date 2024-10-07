package resourcepath

import (
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
)

// composer#taskinstance#DAGID#RUNID#TASKID-MAPINDEX
func ComposerTaskInstance(ti *model.AirflowTaskInstance) string {
	var detail = ti.TaskId()
	if ti.MapIndex() != "-1" {
		detail += "+" + ti.MapIndex()
	}
	return "Cloud Composer#Task Instance#" + strings.ToLower(strings.Join([]string{
		ti.DagId(),
		ti.RunId(),
		detail,
	}, "#"))
}

// composer#airflow-worker#HOST
func ComposerAirflowWorker(wo *model.AirflowWorker) string {
	return "Cloud Composer#Airflow Worker#cluster-scope#" + strings.ToLower(strings.Join([]string{
		wo.Host(),
	}, "#"))
}

func DagFileProcessorStats(stats *model.DagFileProcessorStats) string {
	return "Cloud Composer#Dag File Processor Stats#cluster-scope#" + strings.ToLower(strings.Join([]string{
		stats.DagFilePath(),
	}, "#"))
}
