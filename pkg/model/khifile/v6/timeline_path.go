package khifilev6

import (
	"sync"

	pb "github.com/GoogleCloudPlatform/khi/pkg/generated/khifile/v6"
)

// TimelinePath represents a single node in the timeline hierarchy tree.
// Its identity is uniquely managed by TimelinePathPool, so pointer equality (==)
// guarantees logical equality.
type TimelinePath struct {
	// ID is the unique identifier for this timeline path.
	ID uint32
	// Parent is a pointer to the parent path. It is nil for root paths.
	Parent *TimelinePath
	// Name is a reference to the interned string name of this path segment.
	Name *InternStringRef
	// Type is the style definition for this timeline path segment.
	Type *pb.TimelineType
}

// PathSegment represents a single level of a timeline path to be appended or retrieved.
type PathSegment struct {
	// Name is the string representation of the segment.
	Name string
	// Type is the style definition for this segment.
	Type *pb.TimelineType
}

// timelinePathKey is used internally as a composite key to deduplicate TimelinePath instances.
type timelinePathKey struct {
	// parentID is the ID of the parent TimelinePath, or 0 for root paths.
	parentID uint32
	// nameID is the interned string ID.
	nameID uint32
	// typeID is the TimelineType ID.
	typeID uint32
}

// TimelinePathPool guarantees the uniqueness of TimelinePath instances.
// It is safe for concurrent use by multiple goroutines.
type TimelinePathPool struct {
	// idGen is used to generate unique IDs for new TimelinePath instances.
	idGen *IDGenerator
	// stringPool is used to intern string names of path segments to StringRefs.
	stringPool *InternPool
	// paths is a concurrent map caching timelinePathKey to *TimelinePath.
	paths sync.Map // map[timelinePathKey]*TimelinePath
}

// NewTimelinePathPool creates a new pool for deduplicating TimelinePath instances.
func NewTimelinePathPool(idGen *IDGenerator, sp *InternPool) *TimelinePathPool {
	return &TimelinePathPool{
		idGen:      idGen,
		stringPool: sp,
	}
}

// Get retrieves or creates a multi-level TimelinePath starting from the given parent.
// If parent is nil, it starts from the root.
// This method provides a convenient way to build deep paths in a single call.
func (p *TimelinePathPool) Get(parent *TimelinePath, segments ...PathSegment) *TimelinePath {
	current := parent
	for _, seg := range segments {
		current = p.getOrCreateSingle(current, seg.Name, seg.Type)
	}
	return current
}

// getOrCreateSingle handles the thread-safe retrieval or creation of a single TimelinePath segment.
func (p *TimelinePathPool) getOrCreateSingle(parent *TimelinePath, name string, t *pb.TimelineType) *TimelinePath {
	nameRef := p.stringPool.InternString(name)

	if t == nil || t.Id == nil {
		panic("TimelineType and its ID must not be nil")
	}

	var parentID uint32
	if parent != nil {
		parentID = parent.ID
	}

	key := timelinePathKey{
		parentID: parentID,
		nameID:   nameRef.id,
		typeID:   *t.Id,
	}

	// 1. Fast Path: Load from cache.
	if val, ok := p.paths.Load(key); ok {
		return val.(*TimelinePath)
	}

	// 2. Create new path instance with a fresh ID.
	newPath := &TimelinePath{
		ID:     p.idGen.New(IDTimelinePath),
		Parent: parent,
		Name:   nameRef,
		Type:   t,
	}

	// 3. Atomic store or retrieve if another goroutine won the race.
	actual, _ := p.paths.LoadOrStore(key, newPath)
	return actual.(*TimelinePath)
}
