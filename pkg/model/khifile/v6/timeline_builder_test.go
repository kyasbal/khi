package khifilev6

import (
	"testing"
	"time"

	pb "github.com/GoogleCloudPlatform/khi/pkg/generated/khifile/v6"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimelineBuilder_ToProto_SortsRevisionsByTime(t *testing.T) {
	builder := &TimelineBuilder{}

	now := time.Now()
	t1 := timestamppb.New(now.Add(-2 * time.Hour))
	t2 := timestamppb.New(now.Add(-1 * time.Hour))
	t3 := timestamppb.New(now)

	id1 := uint32(1)
	id2 := uint32(2)
	id3 := uint32(3)

	// Add out of chronological order
	builder.AddRevision(&pb.Revision{LogId: &id2, ChangedTime: t2})
	builder.AddRevision(&pb.Revision{LogId: &id3, ChangedTime: t3})
	builder.AddRevision(&pb.Revision{LogId: &id1, ChangedTime: t1})

	got := builder.ToProto()

	wantRevisions := []*pb.Revision{
		{LogId: &id1, ChangedTime: t1},
		{LogId: &id2, ChangedTime: t2},
		{LogId: &id3, ChangedTime: t3},
	}

	if diff := cmp.Diff(wantRevisions, got.Revisions, protocmp.Transform()); diff != "" {
		t.Errorf("Revisions not sorted chronologically (-want +got):\n%s", diff)
	}
}
