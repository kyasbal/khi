package inspection

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/inspectiondata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"
)

type PrepareInspectionServerFunc = func(inspectionServer *InspectionTaskServer) error

type InspectionType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Priority    int    `json:"-"`
}

type FeatureListItem struct {
	Id          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type InspectionDryRunResult struct {
	Metadata interface{} `json:"metadata"`
}

type InspectionRunResult struct {
	Metadata    interface{}
	ResultStore inspectiondata.Store
}

// InspectionTaskServer manages tasks and provides apis to get task related information in JSON convertible type.
type InspectionTaskServer struct {
	// rootTaskSet is the set of the all definitions in KHI.
	rootTaskSet *task.DefinitionSet
	// inspectionTypes are kinds of tasks. Users will select this at first to filter togglable feature tasks.
	inspectionTypes []*InspectionType
	// tasks are generated tasks
	tasks map[string]*InspectionRunner
}

func NewServer() (*InspectionTaskServer, error) {
	ns, err := task.NewSet([]task.Definition{})
	if err != nil {
		return nil, err
	}
	return &InspectionTaskServer{
		rootTaskSet:     ns,
		inspectionTypes: make([]*InspectionType, 0),
		tasks:           map[string]*InspectionRunner{},
	}, nil
}

// AddInspectionType register a inspection type.
func (s *InspectionTaskServer) AddInspectionType(newInspectionType InspectionType) error {
	if strings.Contains(newInspectionType.Id, "/") {
		return fmt.Errorf("inspection type must not contain /")
	}
	idMap := map[string]interface{}{}
	for _, inspectionType := range s.inspectionTypes {
		idMap[inspectionType.Id] = struct{}{}
	}
	if _, exist := idMap[newInspectionType.Id]; exist {
		return fmt.Errorf("inspection type id:%s is duplicated. InspectionType ID must be unique", newInspectionType.Id)
	}
	inspectionTypesCandidate := append(s.inspectionTypes, &newInspectionType)
	slices.SortFunc(inspectionTypesCandidate, func(a *InspectionType, b *InspectionType) bool {
		return a.Priority > b.Priority
	})
	s.inspectionTypes = inspectionTypesCandidate
	return nil
}

// AddTaskDefinition register a task definition usable for the inspection tasks
func (s *InspectionTaskServer) AddTaskDefinition(taskDefinition task.Definition) error {
	return s.rootTaskSet.Add(taskDefinition)
}

// CreateInspection generates an inspection and returns inspection ID
func (s *InspectionTaskServer) CreateInspection(inspectionType string) (string, error) {
	inspectionTask := NewInspectionRunner(s)
	err := inspectionTask.SetInspectionType(inspectionType)
	if err != nil {
		return "", err
	}
	s.tasks[inspectionTask.ID] = inspectionTask
	return inspectionTask.ID, nil
}

// Inspection returns an instance of an Inspection queried with given inspection ID.
func (s *InspectionTaskServer) GetTask(taskId string) *InspectionRunner {
	return s.tasks[taskId]
}

func (s *InspectionTaskServer) GetAllInspectionTypes() []*InspectionType {
	return append([]*InspectionType{}, s.inspectionTypes...)
}

func (s *InspectionTaskServer) GetInspectionType(inspectionTypeId string) *InspectionType {
	for _, registeredType := range s.inspectionTypes {
		if registeredType.Id == inspectionTypeId {
			return registeredType
		}
	}
	return nil
}

func (s *InspectionTaskServer) GetAllRunners() []*InspectionRunner {
	inspections := []*InspectionRunner{}
	for _, value := range s.tasks {
		inspections = append(inspections, value)
	}
	return inspections
}

// GetAllRegisteredTasks returns a cloned list of all definitions registered in this server.
func (s *InspectionTaskServer) GetAllRegisteredTasks() []task.Definition {
	return s.rootTaskSet.GetAll()
}
