package metadata

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"

// TODO: avoid circular dependency and use namespace in the flag name
const LabelKeyIncludedInRunResultFlag = "metadata/include-in-run-result"
const LabelKeyIncludedInDryRunResultFlag = "metadata/include-in-dry-run-result"
const LabelKeyIncludedInTaskListFlag = "metadata/include-in-tasklist"
const LabelKeyIncludedInResultBinaryFlag = "metadata/include-in-result-binary"

func IncludeInRunResult() task.LabelOpt {
	return task.WithLabel(LabelKeyIncludedInRunResultFlag, true)
}

func IncludeInDryRunResult() task.LabelOpt {
	return task.WithLabel(LabelKeyIncludedInDryRunResultFlag, true)
}

func IncludeInTaskList() task.LabelOpt {
	return task.WithLabel(LabelKeyIncludedInTaskListFlag, true)
}

func IncludeInResultBinary() task.LabelOpt {
	return task.WithLabel(LabelKeyIncludedInResultBinaryFlag, true)
}
