package resourcepath

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
)

// composer#taskinstance#DAGID#RUNID#TASKID-MAPINDEX
func ComposerTaskInstance(ti *model.AirflowTaskInstance) ResourcePath {
	var detail = ti.TaskId()
	if ti.MapIndex() != "-1" {
		detail += "+" + ti.MapIndex()
	}
	return SubresourceLayerGeneralItem("Cloud Composer", "Task Instance", ti.DagId(), ti.RunId(), detail)
}

// composer#airflow-worker#HOST
func ComposerAirflowWorker(wo *model.AirflowWorker) ResourcePath {
	return NameLayerGeneralItem("Cloud Composer", "Airflow Worker", "cluster-scope", wo.Host())
}

func DagFileProcessorStats(stats *model.DagFileProcessorStats) ResourcePath {
	return NameLayerGeneralItem("Cloud Composer", "Dag File Processor Stats", "cluster-scope", stats.DagFilePath())
}
