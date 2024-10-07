package log

import (
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

// Empty commmon log field extractor. Just for testing purpose.
type UnreachableCommonFieldExtractor struct{}

var _ CommonLogFieldExtractor = (*UnreachableCommonFieldExtractor)(nil)

// LogBody implements CommonLogFieldExtractor.
func (u *UnreachableCommonFieldExtractor) LogBody(l *LogEntity) string {
	panic("unimplemented")
}

// DisplayID implements CommonLogFieldExtractor.
func (UnreachableCommonFieldExtractor) DisplayID(l *LogEntity) string {
	panic("unimplemented")
}

func (UnreachableCommonFieldExtractor) ID(l *LogEntity) string {
	panic("Unreachable")
}

func (UnreachableCommonFieldExtractor) Timestamp(l *LogEntity) time.Time {
	panic("Unreachable")
}

func (UnreachableCommonFieldExtractor) MainMessage(l *LogEntity) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (UnreachableCommonFieldExtractor) Severity(l *LogEntity) (enum.Severity, error) {
	return enum.SeverityUnknown, fmt.Errorf("not implemented")
}
