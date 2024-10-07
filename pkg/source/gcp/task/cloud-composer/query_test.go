package composer_task

import (
	"context"
	"fmt"
	"testing"

	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

func TestCreateGeneratorCreatesComposerQuery(t *testing.T) {
	ctx := context.Background()
	projectId := "test-project"
	environmentName := "test-environment"
	vs := &task.VariableSet{}
	vs.Set(gcp_task.InputProjectIdVariableName, projectId)
	vs.Set(InputComposerEnvironmentVariableName, environmentName)

	// resource.type="cloud_composer_environment"
	// resource.labels.environment_name="test-environment"
	// log_name=projects/test-project/logs/airflow-scheduler
	expected := fmt.Sprintf(`resource.type="cloud_composer_environment"
resource.labels.environment_name="test-environment"
log_name=projects/%s/logs/airflow-scheduler`, projectId)

	taskMode := 0                                     // any int is fine
	generator := createGenerator("airflow-scheduler") // sample: airflow-scheduler
	actual, err := generator(ctx, taskMode, vs)
	if err != nil {
		t.Fatalf("GenerateQuery: %v", err)
	}
	if len(actual) != 1 {
		t.Errorf("Unexpected query count %d", len(actual))
	}
	if actual[0] != expected {
		t.Errorf("GenerateQuery: expected %q, got %q", expected, actual)
	}
}
