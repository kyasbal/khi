package query

import (
	"context"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata/query"
	inspection_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	gcp_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"

	api_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/gcp/api"
	inspection_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/inspection"
	task_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/task"
)

type TimestampMockCommonFieldExtractor struct {
	Time time.Time
}

// LogBody implements log.CommonLogFieldExtractor.
func (t *TimestampMockCommonFieldExtractor) LogBody(l *log.LogEntity) string {
	panic("unimplemented")
}

// DisplayID implements log.CommonLogFieldExtractor.
func (*TimestampMockCommonFieldExtractor) DisplayID(l *log.LogEntity) string {
	panic("unimplemented")
}

// Severity implements log.CommonLogFieldExtractor.
func (*TimestampMockCommonFieldExtractor) Severity(l *log.LogEntity) (enum.Severity, error) {
	panic("unimplemented")
}

// ID implements log.CommonLogFieldExtractor.
func (*TimestampMockCommonFieldExtractor) ID(l *log.LogEntity) string {
	panic("unimplemented")
}

// MainMessage implements log.CommonLogFieldExtractor.
func (*TimestampMockCommonFieldExtractor) MainMessage(l *log.LogEntity) (string, error) {
	panic("unimplemented")
}

// Timestamp implements log.CommonLogFieldExtractor.
func (t *TimestampMockCommonFieldExtractor) Timestamp(l *log.LogEntity) time.Time {
	return t.Time
}

var _ log.CommonLogFieldExtractor = (*TimestampMockCommonFieldExtractor)(nil)

func TestNewQueryGeneratorTask(t *testing.T) {
	qg := NewQueryGeneratorTask("foo", "query-name-foo", enum.LogTypeAudit, []string{}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
		return []string{"query-foo"}, nil
	})
	returned := false
	v, err := inspection_test.RunInspectionTaskGraph(qg, map[string]any{},
		api_test.CreateMockAPIClientTask(&api_test.MockApiClient{
			ListLogEntriesFunc: func(ctx context.Context, projectId, filter string, logSink chan any) error {
				if returned {
					close(logSink)
					return nil
				}
				returned = true
				log1 := `test: log1
timestamp: 2023-10-01T12:30:00Z`
				log2 := `test: log2
timestamp: 2023-10-01T12:31:00Z`
				l1Parsed := map[string]string{}
				l2Parsed := map[string]string{}
				err := yaml.Unmarshal([]byte(log1), &l1Parsed)
				if err != nil {
					t.Fatalf("failed to unmarshal log1: %v", err)
				}
				yaml.Unmarshal([]byte(log2), &l2Parsed)
				if err != nil {
					t.Fatalf("failed to unmarshal log2: %v", err)
				}
				logSink <- l1Parsed
				logSink <- l2Parsed
				close(logSink)
				return nil
			},
		}),
		task_test.MockProcessorTaskFromTaskId(inspection_task.ReaderFactoryGeneratorTaskId, structure.NewReaderFactory(&structuredatastore.OnMemoryStructureDataStore{})),
		task_test.MockProcessorTaskFromTaskId(gcp_task.InputProjectIdTask.ID().String(), "project-foo"),
		task_test.MockProcessorTaskFromTaskId(gcp_task.InputStartTimeTask.ID().String(), time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)),
		task_test.MockProcessorTaskFromTaskId(gcp_task.InputEndTimeTask.ID().String(), time.Date(2023, 10, 1, 13, 0, 0, 0, time.UTC)),
		task_test.MockProcessorTaskFromTaskId(gcp_task.InputDurationTask.ID().String(), time.Hour),
	)
	if err != nil {
		t.Errorf("unexpected error\n%v", err)
	}
	m, err := inspection_task.GetMetadataSetFromVariable(v)
	if err != nil {
		t.Errorf("unexpected error\n%v", err)
	}
	q := m.LoadOrStore(query.QueryMetadataKey, &query.QueryMetadataFactory{})
	if diff := cmp.Diff(q, &query.QueryMetadata{
		Queries: []*query.QueryItem{
			{
				Id:    "foo",
				Name:  "query-name-foo",
				Query: "query-foo\ntimestamp >= \"2023-10-01T12:00:00+0000\"\ntimestamp <= \"2023-10-01T13:00:00+0000\"",
			},
		},
	}, cmpopts.IgnoreUnexported(query.QueryMetadata{})); diff != "" {
		t.Errorf("Query metadata is containing non expected value\n%s", diff)
	}

	result, err := v.Get(qg.ID().ReferenceId().String())
	if err != nil {
		t.Errorf("unexpected error\n%v", err)
	}

	if diff := cmp.Diff(result, []*log.LogEntity{
		{LogType: enum.LogTypeAudit},
		{LogType: enum.LogTypeAudit},
	}, cmpopts.IgnoreUnexported(log.LogEntity{}), cmpopts.IgnoreFields(log.LogEntity{}, "Fields")); diff != "" {
		t.Errorf("queried logs are not matching expected value\n%s", diff)
	}
}

func TestNewQueryGeneratorWithDryRunTask(t *testing.T) {
	qg := NewQueryGeneratorTask("foo", "query-name-foo", enum.LogTypeAudit, []string{}, func(ctx context.Context, i int, vs *task.VariableSet) ([]string, error) {
		return []string{"query-foo"}, nil
	})
	v, err := inspection_test.DryRunInspectionTaskGraph(qg, map[string]any{}, api_test.CreateMockAPIClientTask(&api_test.MockApiClient{
		ListLogEntriesFunc: func(ctx context.Context, projectId, filter string, logSink chan any) error {
			t.Errorf("list log api shouldn't be called in dry run mode")
			close(logSink)
			return nil
		},
	}),
		task_test.MockProcessorTaskFromTaskId(inspection_task.ReaderFactoryGeneratorTaskId, structure.NewReaderFactory(&structuredatastore.OnMemoryStructureDataStore{})),
		task_test.MockProcessorTaskFromTaskId(gcp_task.InputProjectIdTask.ID().String(), "project-foo"),
		task_test.MockProcessorTaskFromTaskId(gcp_task.InputStartTimeTask.ID().String(), time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)),
		task_test.MockProcessorTaskFromTaskId(gcp_task.InputEndTimeTask.ID().String(), time.Date(2023, 10, 1, 13, 0, 0, 0, time.UTC)),
		task_test.MockProcessorTaskFromTaskId(gcp_task.InputDurationTask.ID().String(), time.Hour),
	)
	if err != nil {
		t.Errorf("unexpected error\n%v", err)
	}
	m, err := inspection_task.GetMetadataSetFromVariable(v)
	if err != nil {
		t.Errorf("unexpected error\n%v", err)
	}
	q := m.LoadOrStore(query.QueryMetadataKey, &query.QueryMetadataFactory{})
	if diff := cmp.Diff(q, &query.QueryMetadata{
		Queries: []*query.QueryItem{
			{
				Id:    "foo",
				Name:  "query-name-foo",
				Query: "query-foo\ntimestamp >= \"2023-10-01T12:00:00+0000\"\ntimestamp <= \"2023-10-01T13:00:00+0000\"",
			},
		},
	}, cmpopts.IgnoreUnexported(query.QueryMetadata{})); diff != "" {
		t.Errorf("Query metadata is containing non expected value\n%s", diff)
	}

	result, err := v.Get(qg.ID().ReferenceId().String())
	if err != nil {
		t.Errorf("unexpected error\n%v", err)
	}

	if diff := cmp.Diff(result, []*log.LogEntity{}, cmpopts.IgnoreUnexported(log.LogEntity{})); diff != "" {
		t.Errorf("queried logs are not matching expected value\n%s", diff)
	}
}
