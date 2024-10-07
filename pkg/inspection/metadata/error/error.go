package error

import (
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

const ErrorMessageSetMetadataKey = "error"

type ErrorMessage struct {
	ErrorId int    `json:"errorId"`
	Message string `json:"message"`
	Link    string `json:"link"`
}

// ErrorMessageSet is a metadata type containing errors exposed to frontend.
type ErrorMessageSet struct {
	ErrorMessages []*ErrorMessage `json:"errorMessages"`
}

// Labels implements metadata.Metadata.
func (e *ErrorMessageSet) Labels() *task.LabelSet {
	return task.NewLabelSet(metadata.IncludeInRunResult(), metadata.IncludeInTaskList())
}

// ToSerializable implements metadata.Metadata.
func (e *ErrorMessageSet) ToSerializable() interface{} {
	return e
}

var _ metadata.Metadata = (*ErrorMessageSet)(nil)

// AddErrorMessage stores a new ErrorMessage. Duplicated error message will be ignored.
func (e *ErrorMessageSet) AddErrorMessage(newError *ErrorMessage) {
	for _, msg := range e.ErrorMessages {
		if msg.ErrorId == newError.ErrorId {
			return // Skip adding duplicated error
		}
	}
	e.ErrorMessages = append(e.ErrorMessages, newError)
}

func NewPermissionErrorMessage(projectId string) *ErrorMessage {
	return &ErrorMessage{
		ErrorId: 0,
		Message: fmt.Sprintf("Permission error to read logs from project `%s`", projectId),
		Link:    "https://g3doc.corp.google.com/company/gfw/support/cloud/systems/khi/troubleshooting.md?cl=head",
	}
}

func NewNotFoundErrorMessage(projectId string) *ErrorMessage {
	return &ErrorMessage{
		ErrorId: 1,
		Message: fmt.Sprintf("Project `%s` not found", projectId),
		Link:    "https://g3doc.corp.google.com/company/gfw/support/cloud/systems/khi/troubleshooting.md?cl=head",
	}
}

func NewUnauthorizedErrorMessage() *ErrorMessage {
	return &ErrorMessage{
		ErrorId: 2,
		Message: "Access token is not authorized. (Token expired?)",
		Link:    "https://g3doc.corp.google.com/company/gfw/support/cloud/systems/khi/troubleshooting.md?cl=head",
	}
}

var _ metadata.MetadataFactory = (*ErrorMessageSetFactory)(nil)

type ErrorMessageSetFactory struct {
}

// Instanciate implements metadata.MetadataFactory.
func (e *ErrorMessageSetFactory) Instanciate() metadata.Metadata {
	return &ErrorMessageSet{[]*ErrorMessage{}}
}
