package testlog

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/adapter"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log/structure/structuredatastore"
	"gopkg.in/yaml.v3"
)

type TestLogOpt = func(original *yaml.Node) (*yaml.Node, error)

// TestLog is a type to generate mock log data effectively for test.
type TestLog struct {
	opts []TestLogOpt
}

func New(opts ...TestLogOpt) *TestLog {
	return &TestLog{
		opts: opts,
	}
}

// With instantiate a new TestLog with the given additional options and the current options.
func (b *TestLog) With(additionalOpts ...TestLogOpt) *TestLog {
	opts := []TestLogOpt{}
	opts = append(opts, b.opts...)
	opts = append(opts, additionalOpts...)
	return &TestLog{
		opts: opts,
	}
}

func (b *TestLog) BuildReader() (*structure.Reader, error) {
	var node *yaml.Node
	var err error
	for _, opt := range b.opts {
		node, err = opt(node)
		if err != nil {
			return nil, err
		}
	}
	sd, err := structuredata.DataFromYamlNode(node)
	if err != nil {
		return nil, err
	}
	directStore := adapter.Direct(sd)
	return directStore.GetReaderBackedByStore(&structuredatastore.OnMemoryStructureDataStore{})
}

func (b *TestLog) MustBuildYamlString() string {
	reader, err := b.BuildReader()
	if err != nil {
		panic(err)
	}
	yamlStr, err := reader.ToYaml("")
	if err != nil {
		panic(err)
	}
	return yamlStr
}

func (b *TestLog) MustBuildLogEntity(le log.CommonLogFieldExtractor) *log.LogEntity {
	reader, err := b.BuildReader()
	if err != nil {
		panic(err)
	}
	return log.NewLogEntity(reader, le)
}
