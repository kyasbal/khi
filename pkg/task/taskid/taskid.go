package taskid

import (
	"slices"
	"strings"
)

// TaskImplementationId is a unique value associated with each task implementation.
// TaskReferenceId is an ID that can be used to refer tasks in the dependencies and does not
// have a task suffix.

// For example, "foo.bar" can be used as both a TaskImplementationId and a
// TaskReferenceId, but "foo.bar#qux" cannot be used as a TaskReferenceId because it
// has a suffix.

type TaskReferenceId struct {
	id string
}
type TaskImplementationId struct {
	referenceId        string
	implementationHash string
}

func (r TaskReferenceId) String() string {
	return r.id
}

func (i TaskImplementationId) String() string {
	if i.implementationHash == "" {
		return i.referenceId
	}
	return i.referenceId + "#" + i.implementationHash
}

func (i TaskImplementationId) ImplementationHash() string {
	return i.implementationHash
}

func (i TaskImplementationId) ReferenceId() TaskReferenceId {
	return NewTaskReference(i.referenceId)
}

func NewTaskReference(taskId string) TaskReferenceId {
	return TaskReferenceId{
		id: taskId,
	}
}

func NewTaskImplementationId(taskId string) TaskImplementationId {
	if strings.Count(taskId, "#") > 0 {
		splitted := strings.Split(taskId, "#")
		return TaskImplementationId{
			referenceId:        splitted[0],
			implementationHash: splitted[1],
		}
	} else {
		return TaskImplementationId{
			referenceId:        taskId,
			implementationHash: "",
		}
	}
}

func (i TaskImplementationId) Match(reference TaskReferenceId) bool {
	return i.referenceId == reference.id
}

func DedupeReferenceIds(references []TaskReferenceId) []TaskReferenceId {
	found := map[string]struct{}{}
	for _, elem := range references {
		found[elem.id] = struct{}{}
	}
	result := []TaskReferenceId{}
	for id := range found {
		result = append(result, TaskReferenceId{
			id: id,
		})
	}
	slices.SortFunc(result, func(a, b TaskReferenceId) int { return strings.Compare(a.id, b.id) })
	return result
}
