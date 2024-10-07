package history

import (
	"time"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/log"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/binarychunk"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"
)

// Staging data types contains large string data directly on it.
// These types are just used in the arguments of ChangeSet.
// These will be converted to the corresponding serializable types with some binary data stored in binarychunk.Builder.

type StagingResourceRevision struct {
	Verb      enum.RevisionVerb
	Body      string
	Requestor string
	Partial   bool
	// If this resource existence is inferred from another logs later.
	Inferred   bool
	ChangeTime time.Time
	State      enum.RevisionState
}

func (r *StagingResourceRevision) commit(binaryBuilder *binarychunk.Builder, l *log.LogEntity) (*ResourceRevision, error) {
	bodyRef, err := binaryBuilder.Write([]byte(r.Body))
	if err != nil {
		return nil, err
	}
	requestorRef, err := binaryBuilder.Write([]byte(r.Requestor))
	if err != nil {
		return nil, err
	}
	logId := l.ID()
	if r.Inferred {
		logId = ""
	}
	return &ResourceRevision{
		Log:        logId,
		Requestor:  requestorRef,
		Verb:       r.Verb,
		Body:       bodyRef,
		Partial:    r.Partial,
		ChangeTime: r.ChangeTime,
		State:      r.State,
	}, nil
}
