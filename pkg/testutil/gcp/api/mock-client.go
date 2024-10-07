package api_test

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/progress"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/api"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

type MockApiClient struct {
	GetClusterNamesFunc func(ctx context.Context, projectId string) ([]string, error)
	ListLogEntriesFunc  func(ctx context.Context, projectId string, filter string, logSink chan any) error
}

// GetAnthosAWSClusterNames implements api.GCPClient.
func (m *MockApiClient) GetAnthosAWSClusterNames(ctx context.Context, projectId string) ([]string, error) {
	if m.GetClusterNamesFunc == nil {
		return []string{"aws-cluster-foo", "aws-cluster-bar"}, nil
	}
	return m.GetClusterNamesFunc(ctx, projectId)
}

// GetAnthosAzureClusterNames implements api.GCPClient.
func (m *MockApiClient) GetAnthosAzureClusterNames(ctx context.Context, projectId string) ([]string, error) {
	if m.GetClusterNamesFunc == nil {
		return []string{"azure-cluster-foo", "azure-cluster-bar"}, nil
	}
	return m.GetClusterNamesFunc(ctx, projectId)
}

// GetAnthosOnBaremetalClusterNames implements api.GCPClient.
func (m *MockApiClient) GetAnthosOnBaremetalClusterNames(ctx context.Context, projectId string) ([]string, error) {
	if m.GetClusterNamesFunc == nil {
		return []string{"baremetal-cluster-foo", "baremetal-cluster-bar"}, nil
	}
	return m.GetClusterNamesFunc(ctx, projectId)
}

// GetAnthosOnVMWareClusterNames implements api.GCPClient.
func (m *MockApiClient) GetAnthosOnVMWareClusterNames(ctx context.Context, projectId string) ([]string, error) {
	if m.GetClusterNamesFunc == nil {
		return []string{"vmware-cluster-foo", "vmware-cluster-bar"}, nil
	}
	return m.GetClusterNamesFunc(ctx, projectId)
}

// GetClusterNames implements api.GCPClient.
func (m *MockApiClient) GetClusterNames(ctx context.Context, projectId string) ([]string, error) {
	if m.GetClusterNamesFunc == nil {
		return []string{"gke-cluster-foo", "gke-cluster-bar"}, nil
	}
	return m.GetClusterNamesFunc(ctx, projectId)
}

func (m *MockApiClient) GetComposerEnvironmentNames(ctx context.Context, projectId string, location string) ([]string, error) {
	// GetClusterNamesFunc is not for Composer environment? Yes, but it's fine since it is a mock! :D
	if m.GetClusterNamesFunc == nil {
		return []string{"composer-environment-foo", "composer-environment-bar"}, nil
	}
	return m.GetClusterNamesFunc(ctx, projectId)
}

// ListLogEntries implements api.GCPClient.
func (m *MockApiClient) ListLogEntries(ctx context.Context, projectId string, filter string, logSink chan any) error {
	if m.ListLogEntriesFunc == nil {
		close(logSink)
		return nil
	}
	return m.ListLogEntriesFunc(ctx, projectId, filter, logSink)
}

var _ api.GCPClient = (*MockApiClient)(nil)

func CreateMockAPIClientTask(client *MockApiClient) task.Definition {
	return inspection_task.NewInspectionProducer(gcp_task.GCPApiClientTaskId, func(ctx context.Context, taskMode int, progress *progress.TaskProgress) (any, error) {
		return client, nil
	})
}
