package history

import (
	"context"
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/worker"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/ioconfig"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
	gcp_log "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/source/gcp/log"
	log_test "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/testutil/testlog"
	"github.com/google/go-cmp/cmp"
)

func TestRecordLogSummary(t *testing.T) {
	log := log_test.MockLogWithId("foo")
	cs := NewChangeSet(log)
	cs.RecordLogSummary("bar")
	if cs.logSummaryRewrite != "bar" {
		t.Errorf("logSummaryRewrite is not rewritten to the expected value")
	}
}

func TestRecordLogSeverity(t *testing.T) {
	log := log_test.MockLogWithId("foo")
	cs := NewChangeSet(log)
	cs.RecordLogSeverity(enum.SeverityWarning)
	if cs.logSeverityRewrite != enum.SeverityWarning {
		t.Errorf("logSeverityRewrite is not rewritten to the expected value")
	}
}

func TestRecordEvents(t *testing.T) {
	log := log_test.MockLogWithId("foo")
	cs := NewChangeSet(log)
	cs.RecordEvent("A#B")
	cs.RecordEvent("A#C", RewriteRelationship(enum.RelationshipOperation))
	if diff := cmp.Diff(cs.events, map[string][]*ResourceEvent{
		"A#B": {{Log: "foo"}},
		"A#C": {{Log: "foo"}},
	}); diff != "" {
		t.Errorf("RecordEvent didn't modify ChangeSet as expected\n%s", diff)
	}
	if diff := cmp.Diff(cs.resourceOpts, map[string][]ResourceOpt{
		"A#C": {&rewriteRelationshipImpl{relationship: enum.RelationshipOperation}},
	}, cmp.AllowUnexported(rewriteRelationshipImpl{})); diff != "" {
		t.Errorf("RecordEvent didn't modify resourceOpts in ChangeSet as expected\n%s", diff)
	}
}

func TestRecordRevisions(t *testing.T) {
	log := log_test.MockLogWithId("foo")
	cs := NewChangeSet(log)
	cs.RecordRevision("A#B", &StagingResourceRevision{
		Inferred: true,
	}, RewriteRelationship(enum.RelationshipContainer))
	cs.RecordRevision("A#B", &StagingResourceRevision{})
	cs.RecordRevision("A#C", &StagingResourceRevision{})
	if diff := cmp.Diff(cs.revisions, map[string][]*StagingResourceRevision{
		"A#B": {{Inferred: true}, {}},
		"A#C": {{}},
	}); diff != "" {
		t.Errorf("RecordRevision didn't modify ChangeSet as expected\n%s", diff)
	}

	if diff := cmp.Diff(cs.annotations, []LogAnnotation{
		&ResourceReferenceAnnotation{Path: "A#B"},
		&ResourceReferenceAnnotation{Path: "A#C"},
	}); diff != "" {
		t.Errorf("RecordRevision didn't modify log annotations in ChangeSet as expected\n%s", diff)
	}
	if diff := cmp.Diff(cs.resourceOpts, map[string][]ResourceOpt{
		"A#B": {&rewriteRelationshipImpl{relationship: enum.RelationshipContainer}},
	}, cmp.AllowUnexported(rewriteRelationshipImpl{})); diff != "" {
		t.Errorf("RecordRevision didn't modify resourceOpts in ChangeSet as expected\n%s", diff)
	}
}

func TestChangesetFlushIsThreadSafe(t *testing.T) {
	groupCount := 100
	logCountPerGroup := 100
	builder := NewBuilder(&ioconfig.IOConfig{})
	lt := testlog.New(testlog.BaseYaml(""))
	l := [][]*log.LogEntity{}
	allLogs := []*log.LogEntity{}
	for i := 0; i < groupCount; i++ {
		l = append(l, make([]*log.LogEntity, 0))
	}
	for li := 0; li < logCountPerGroup; li++ {
		for i := 0; i < groupCount; i++ {
			hour := i / 3600
			minute := (i - hour*3600) / 60
			seconds := (i - hour*3600 - minute*60) % 60
			l[i] = append(l[i], lt.With(
				testlog.StringField("insertId", fmt.Sprintf("id-group%d-%d", i, li)),
				testlog.StringField("timestamp", fmt.Sprintf("2024-01-01T%02d:%02d:%02dZ", hour, minute, seconds)),
			).MustBuildLogEntity(gcp_log.GCPCommonFieldExtractor{}))
		}
	}
	for _, group := range l {
		allLogs = append(allLogs, group...)
	}
	err := builder.PrepareParseLogs(context.Background(), allLogs, func() {})
	if err != nil {
		t.Fatal(err.Error())
	}
	pool := worker.NewPool(groupCount)
	for i := 0; i < groupCount; i++ {
		currentGroup := l[i]
		groupPath := fmt.Sprintf("grp#%d", i)
		pool.Run(func() {
			for _, l := range currentGroup {
				cs := NewChangeSet(l)
				cs.RecordRevision(groupPath, &StagingResourceRevision{})
				paths, err := cs.FlushToHistory(builder)

				for _, path := range paths {
					tb := builder.GetTimelineBuilder(path)
					tb.Sort()
				}
				if err != nil {
					t.Fatal(err.Error())
				}
			}
		})
	}

	pool.Wait()
	for i := 0; i < groupCount; i++ {
		grpPath := fmt.Sprintf("grp#%d", i)
		tb := builder.GetTimelineBuilder(grpPath)
		if len(tb.timeline.Revisions) != logCountPerGroup {
			t.Errorf("revision count mismatch: expected %d, actual %d", logCountPerGroup, len(tb.timeline.Revisions))
		}
		for li := 0; li < logCountPerGroup; li++ {
			rev := tb.timeline.Revisions[li]
			sl, err := tb.builder.GetLog(rev.Log)
			expectedId := fmt.Sprintf("id-group%d-%d", i, li)
			if err != nil {
				t.Errorf("log %s not found!", rev.Log)
				continue
			}
			if sl.DisplayId != expectedId {
				t.Errorf("log id mismatch: expected %s, actual %s", expectedId, sl.DisplayId)
			}
		}
	}
}
