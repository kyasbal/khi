package log_test

import (
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/adapter"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	gcp_log "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/log"
)

type mockLogFieldExtractor struct {
	id          string
	mainMessage string
	severity    enum.Severity
	time        time.Time
}

// LogBody implements log.CommonLogFieldExtractor.
func (m *mockLogFieldExtractor) LogBody(l *log.LogEntity) string {
	panic("unimplemented")
}

// DisplayID implements log.CommonLogFieldExtractor.
func (*mockLogFieldExtractor) DisplayID(l *log.LogEntity) string {
	panic("unimplemented")
}

// ID implements log.CommonLogFieldExtractor.
func (m *mockLogFieldExtractor) ID(l *log.LogEntity) string {
	return m.id
}

// MainMessage implements log.CommonLogFieldExtractor.
func (m *mockLogFieldExtractor) MainMessage(l *log.LogEntity) (string, error) {
	return m.mainMessage, nil
}

// Severity implements log.CommonLogFieldExtractor.
func (m *mockLogFieldExtractor) Severity(l *log.LogEntity) (enum.Severity, error) {
	return m.severity, nil
}

// Timestamp implements log.CommonLogFieldExtractor.
func (m *mockLogFieldExtractor) Timestamp(l *log.LogEntity) time.Time {
	return m.time
}

var _ log.CommonLogFieldExtractor = (*mockLogFieldExtractor)(nil)

func MustLogEntity(text string) *log.LogEntity {
	readerFactory := structure.NewReaderFactory(&structuredatastore.OnMemoryStructureDataStore{})
	reader, err := readerFactory.NewReader(adapter.Yaml(text))
	if err != nil {
		panic(err)
	}
	yaml := log.NewLogEntity(reader, gcp_log.GCPCommonFieldExtractor{})
	return yaml
}

func MockLogWithId(id string) *log.LogEntity {
	readerFactory := structure.NewReaderFactory(&structuredatastore.OnMemoryStructureDataStore{})
	reader, err := readerFactory.NewReader(adapter.Yaml(fmt.Sprintf("# mock log for %s", id)))
	if err != nil {
		panic(err)
	}
	yaml := log.NewLogEntity(reader, &mockLogFieldExtractor{
		id:          id,
		mainMessage: fmt.Sprintf("# mock log for %s", id),
	})
	if err != nil {
		panic(err)
	}
	return yaml
}
