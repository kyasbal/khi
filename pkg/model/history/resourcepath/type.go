package resourcepath

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"

// ResourcePath contains the path representing location of a timeline in the history.
type ResourcePath struct {
	// Path is the the raw resource path represented in string like `A#B#C`. This means `the C under the B under the A at the root`.
	// KHI uses `#` as the delimiter of resource paths, this is because the root element(API version) can contain `.` or `/`.
	Path string

	// ParentRelationship explains between the location represented with this ResourcePath and its parent.
	// KHI shows various resources in a single history with mixing many types of children. It's not only like child-parent relationship, but also pseudo relationship like node-node's component relationship.
	ParentRelationship enum.ParentRelationship
}
