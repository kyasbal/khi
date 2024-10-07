package history

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/model/binarychunk"

type LogAnnotation interface {
	Priority() int
	Serialize(builder *binarychunk.Builder) (any, error)
}

type ResourceReferenceAnnotation struct {
	Path string
}

type SerializableResourceReferenceAnnotation struct {
	Type string                       `json:"type"`
	Path *binarychunk.BinaryReference `json:"path"`
}

var _ LogAnnotation = (*ResourceReferenceAnnotation)(nil)

func (a *ResourceReferenceAnnotation) Priority() int {
	return 10000
}

func (a *ResourceReferenceAnnotation) Serialize(builder *binarychunk.Builder) (any, error) {
	ref, err := builder.Write([]byte(a.Path))
	if err != nil {
		return nil, err
	}
	return &SerializableResourceReferenceAnnotation{
		Type: "resource_ref",
		Path: ref,
	}, nil
}

func NewResourceReferenceAnnotation(resourcePath string) *ResourceReferenceAnnotation {
	return &ResourceReferenceAnnotation{
		Path: resourcePath,
	}
}
