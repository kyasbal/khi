package api

import (
	"context"
)

type GCPClient interface {
	GetClusterNames(ctx context.Context, projectId string) ([]string, error)
	GetAnthosAWSClusterNames(ctx context.Context, projectId string) ([]string, error)
	GetAnthosAzureClusterNames(ctx context.Context, projectId string) ([]string, error)
	GetAnthosOnBaremetalClusterNames(ctx context.Context, projectId string) ([]string, error)
	GetAnthosOnVMWareClusterNames(ctx context.Context, projectId string) ([]string, error)
	GetComposerEnvironmentNames(ctx context.Context, projectId string, location string) ([]string, error)
	ListLogEntries(ctx context.Context, projectId string, filter string, logSink chan any) error
}

type RefreshableToken interface {
	Refresh() (string, error)
}
