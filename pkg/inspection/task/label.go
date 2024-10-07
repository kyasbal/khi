package task

import common_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"

const (
	InspectionTaskPrefix                 = common_task.KHISystemPrefix + "inspection/"
	LabelKeyInspectionFeatureFlag        = InspectionTaskPrefix + "feature"
	LabelKeyInspectionDefaultFeatureFlag = InspectionTaskPrefix + "default-feature"
	LabelKeyInspectionRequiredFlag       = InspectionTaskPrefix + "required"
	LabelKeyProgressReportable           = InspectionTaskPrefix + "progress-reportable"
	// A []string typed label of Definition. Task registry will filter task units by given inspection type at first.
	LabelKeyInspectionTypes  = InspectionTaskPrefix + "inspection-type"
	LabelKeyFeatureTaskTitle = InspectionTaskPrefix + "feature/title"

	LabelKeyFeatureTaskDescription = InspectionTaskPrefix + "feature/description"

	InspectionMainSubgraphName = InspectionTaskPrefix + "inspection-main"

	TaskModeDryRun = 1
	TaskModeRun    = 2
)

type ProgressReportableTaskLabelOptImpl struct{}

// Write implements task.LabelOpt.
func (i *ProgressReportableTaskLabelOptImpl) Write(label *common_task.LabelSet) {
	label.Set(LabelKeyProgressReportable, true)
}

var _ common_task.LabelOpt = (*ProgressReportableTaskLabelOptImpl)(nil)

// FeatureTaskLabelImpl is an implementation of task.LabelOpt.
// This annotate a task definition to be a feature in inspection.
type FeatureTaskLabelImpl struct {
	title            string
	description      string
	isDefaultFeature bool
}

func (ftl *FeatureTaskLabelImpl) Write(label *common_task.LabelSet) {
	label.Set(LabelKeyInspectionFeatureFlag, true)
	label.Set(LabelKeyFeatureTaskTitle, ftl.title)
	label.Set(LabelKeyFeatureTaskDescription, ftl.description)
	label.Set(LabelKeyInspectionDefaultFeatureFlag, ftl.isDefaultFeature)
}

func (ftl *FeatureTaskLabelImpl) WithDescription(description string) *FeatureTaskLabelImpl {
	ftl.description = description
	return ftl
}

var _ common_task.LabelOpt = (*FeatureTaskLabelImpl)(nil)

func FeatureTaskLabel(title string, description string, isDefaultFeature bool) *FeatureTaskLabelImpl {
	return &FeatureTaskLabelImpl{
		title:            title,
		description:      description,
		isDefaultFeature: isDefaultFeature,
	}
}

type InspectionTaskLabelImpl struct {
	inspectionTypes []string
}

// Write implements task.LabelOpt.
func (itl *InspectionTaskLabelImpl) Write(label *common_task.LabelSet) {
	label.Set(LabelKeyInspectionTypes, itl.inspectionTypes)
}

var _ common_task.LabelOpt = (*InspectionTaskLabelImpl)(nil)

func InspectionTaskLabel(types ...string) *InspectionTaskLabelImpl {
	return &InspectionTaskLabelImpl{
		inspectionTypes: types,
	}
}

type RequriredTaskLabelImpl struct{}

func (r *RequriredTaskLabelImpl) Write(label *common_task.LabelSet) {
	label.Set(LabelKeyInspectionRequiredFlag, true)
}

func NewRequiredTaskLabel() *RequriredTaskLabelImpl {
	return &RequriredTaskLabelImpl{}
}
