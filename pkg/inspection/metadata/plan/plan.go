package plan

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

var InspectionPlanMetadataKey = "plan"

type InspectionPlan struct {
	TaskGraph string `json:"taskGraph"`
}

// Labels implements metadata.Metadata.
func (*InspectionPlan) Labels() *task.LabelSet {
	return task.NewLabelSet(metadata.IncludeInDryRunResult(), metadata.IncludeInRunResult())
}

// ToSerializable implements metadata.Metadata.
func (p *InspectionPlan) ToSerializable() interface{} {
	return p
}

var _ metadata.Metadata = (*InspectionPlan)(nil)

type InspectionPlanMetadataFactory struct{}

// Instanciate implements metadata.MetadataFactory.
func (i *InspectionPlanMetadataFactory) Instanciate() metadata.Metadata {
	return &InspectionPlan{}
}

// InspectionPlanMetadataFactory implements metadata.MetadataFactory
var _ metadata.MetadataFactory = (*InspectionPlanMetadataFactory)(nil)
