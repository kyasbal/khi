package parser_test

import (
	"context"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/history"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/parser"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil"
	log_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/log"
)

// ParseFromYamlLogFile returns the parsed ChangeSet from the yaml log file at the given path with specified parser.
func ParseFromYamlLogFile(testFile string, parser parser.Parser, builder *history.Builder, variables *task.VariableSet) (*history.ChangeSet, error) {
	testutil.InitTestIO()
	yamlStr := testutil.MustReadText(testFile)
	l := log_test.MustLogEntity(yamlStr)
	cs := history.NewChangeSet(l)
	err := parser.Parse(context.Background(), l, cs, builder, variables)
	if err != nil {
		return nil, err
	}
	return cs, nil
}
