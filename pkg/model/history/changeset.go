package history

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

// history.ChangeSet is set of changes applicable to history.
// A parser ingest a log.LogEntry and returns a ChangeSet. ChangeSet contains multiple changes against the history.
// This change is applied atomically, when the parser returns an error, no partial changes would be written.
type ChangeSet struct {
	associatedLog      *log.LogEntity
	revisions          map[string][]*StagingResourceRevision
	events             map[string][]*ResourceEvent
	resourceOpts       map[string][]ResourceOpt
	annotations        []LogAnnotation
	logSummaryRewrite  string
	logSeverityRewrite enum.Severity
	aliases            map[string][]string
	aliasRelationships map[string]enum.ParentRelationShip
}

func NewChangeSet(l *log.LogEntity) *ChangeSet {
	return &ChangeSet{
		associatedLog:      l,
		revisions:          make(map[string][]*StagingResourceRevision),
		events:             make(map[string][]*ResourceEvent),
		resourceOpts:       make(map[string][]ResourceOpt),
		logSummaryRewrite:  "",
		logSeverityRewrite: enum.SeverityUnknown,
		annotations:        []LogAnnotation{},
		aliases:            map[string][]string{},
		aliasRelationships: make(map[string]enum.ParentRelationShip),
	}
}

func (cs *ChangeSet) RecordLogSummary(summary string) {
	cs.logSummaryRewrite = summary
}

func (cs *ChangeSet) RecordLogSeverity(severity enum.Severity) {
	cs.logSeverityRewrite = severity
}

func (cs *ChangeSet) RecordRevision(resourcePath string, revision *StagingResourceRevision, changeSetOpts ...ResourceOpt) {
	if _, exist := cs.revisions[resourcePath]; !exist {
		cs.revisions[resourcePath] = make([]*StagingResourceRevision, 0)
	}
	cs.revisions[resourcePath] = append(cs.revisions[resourcePath], revision)
	if !revision.Inferred {
		cs.annotations = append(cs.annotations, NewResourceReferenceAnnotation(resourcePath))
	}
	for _, opt := range changeSetOpts {
		cs.storeChangeSetOpt(resourcePath, opt)
	}
}

func (cs *ChangeSet) RecordEvent(resourcePath string, changeSetOpts ...ResourceOpt) {
	event := ResourceEvent{
		Log: cs.associatedLog.ID(),
	}
	if _, exist := cs.events[resourcePath]; !exist {
		cs.events[resourcePath] = make([]*ResourceEvent, 0)
	}
	cs.events[resourcePath] = append(cs.events[resourcePath], &event)
	cs.annotations = append(cs.annotations, NewResourceReferenceAnnotation(resourcePath))
	for _, opt := range changeSetOpts {
		cs.storeChangeSetOpt(resourcePath, opt)
	}
}

func (cs *ChangeSet) RecordAliasRelationship(source string, dest string, changeSetOpts ...ResourceOpt) {
	if _, exist := cs.aliases[source]; !exist {
		cs.aliases[source] = make([]string, 0)
	}
	for _, d := range cs.aliases[source] {
		if d == dest {
			return
		}
	}
	cs.aliases[source] = append(cs.aliases[source], dest)
	for _, opt := range changeSetOpts {
		cs.storeChangeSetOpt(dest, opt)
	}
}

func (cs *ChangeSet) storeChangeSetOpt(path string, changeSetOpt ResourceOpt) {
	if _, exist := cs.resourceOpts[path]; !exist {
		cs.resourceOpts[path] = make([]ResourceOpt, 0)
	}
	cs.resourceOpts[path] = append(cs.resourceOpts[path], changeSetOpt)
}

// FlushToHistory export recorded changes to the target history
func (cs *ChangeSet) FlushToHistory(builder *Builder) ([]string, error) {
	changedPaths := []string{}
	for resourcePath, revisions := range cs.revisions {
		tb := builder.GetTimelineBuilder(resourcePath)
		for _, stagingRevision := range revisions {
			revision, err := stagingRevision.commit(builder.binaryChunk, cs.associatedLog)
			if err != nil {
				return nil, err
			}
			tb.AddRevision(revision)
		}
		changedPaths = append(changedPaths, resourcePath)
	}
	for resourcePath, events := range cs.events {
		tb := builder.GetTimelineBuilder(resourcePath)
		for _, event := range events {
			tb.AddEvent(event)
		}
		changedPaths = append(changedPaths, resourcePath)
	}
	if cs.logSummaryRewrite != "" {
		builder.setLogSummary(cs.associatedLog.ID(), cs.logSummaryRewrite)
	}
	if cs.logSeverityRewrite != enum.SeverityUnknown {
		builder.setLogSeverity(cs.associatedLog.ID(), cs.logSeverityRewrite)
	}
	builder.setLogAnnotations(cs.associatedLog.ID(), cs.annotations)
	for source, destinations := range cs.aliases {
		for _, dest := range destinations {
			builder.addTimelineAlias(source, dest, cs.aliasRelationships[dest])
		}
	}
	for resourcePath, changeOpts := range cs.resourceOpts {
		for _, changeOpt := range changeOpts {
			err := changeOpt.Write(builder, resourcePath)
			if err != nil {
				return nil, err
			}
		}
	}
	return common.DedupeStringArray(changedPaths), nil
}
