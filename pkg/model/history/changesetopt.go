package history

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/enum"

type ResourceOpt interface {
	Write(builder *Builder, path string) error
}

type rewriteRelationshipImpl struct {
	relationship enum.ParentRelationShip
}

// Write implements ResourceOpt.
func (i *rewriteRelationshipImpl) Write(builder *Builder, path string) error {
	return builder.rewriteRelationship(path, i.relationship)
}

var _ ResourceOpt = (*rewriteRelationshipImpl)(nil)

// RewriteRelationship returns a ResourceOpt implementation to rewrite the resource relationship type to the specified type.
func RewriteRelationship(relationship enum.ParentRelationShip) ResourceOpt {
	return &rewriteRelationshipImpl{
		relationship: relationship,
	}
}
