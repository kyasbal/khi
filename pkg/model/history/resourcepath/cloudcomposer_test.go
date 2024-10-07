package resourcepath

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model"
)

func TestComposerTaskInstance(t *testing.T) {
	tests := []struct {
		name string
		ti   *model.AirflowTaskInstance
		want string
	}{
		{
			name: "basic",
			ti:   model.NewAirflowTaskInstance("my_dag", "my_task", "my_run", "0", "my_host", "my_status"),
			want: "Cloud Composer#Task Instance#my_dag#my_run#my_task+0",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ComposerTaskInstance(test.ti)
			if got != test.want {
				t.Errorf("ComposerTaskInstance(%v) = %v, want %v", test.ti, got, test.want)
			}
		})
	}
}

func TestComposerAirflowWorker(t *testing.T) {
	tests := []struct {
		name string
		wo   *model.AirflowWorker
		want string
	}{
		{
			name: "basic",
			wo:   model.NewAirflowWorker("my_host"),
			want: "Cloud Composer#Airflow Worker#cluster-scope#my_host",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ComposerAirflowWorker(test.wo)
			if got != test.want {
				t.Errorf("ComposerAirflowWorker(%v) = %v, want %v", test.wo, got, test.want)
			}
		})
	}
}

func TestDagFileProcessorStats(t *testing.T) {
	tests := []struct {
		name  string
		stats *model.DagFileProcessorStats
		want  string
	}{
		{
			name:  "basic",
			stats: model.NewDagFileProcessorStats("my_dag_file_path", "my_dag_file_path", "10", "10"),
			want:  "Cloud Composer#Dag File Processor Stats#cluster-scope#my_dag_file_path",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := DagFileProcessorStats(test.stats)
			if got != test.want {
				t.Errorf("DagFileProcessorStats(%v) = %v, want %v", test.stats, got, test.want)
			}
		})
	}
}
