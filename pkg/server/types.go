package server

import (
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
)

type SerializedMetadata = map[string]any

type ServerStat struct {
	TotalMemoryAvailable int `json:"totalMemoryAvailable"`
}

// GetInspectionTypesResponse is the type of the response for /api/v2/inspection/types
type GetInspectionTypesResponse struct {
	Types []*inspection.InspectionType `json:"types"`
}

// GetInspectionTasksResponse is the type of the response for /api/v2/inspection/tasks
type GetInspectionTasksResponse struct {
	Tasks      map[string]SerializedMetadata `json:"tasks"`
	ServerStat *ServerStat                   `json:"serverStat"`
}

type PostInspectionTaskResponse struct {
	InspectionId string `json:"inspectionId"`
}

type PutInspectionTaskFeatureRequest struct {
	Features []string `json:"features"`
}

type PutInspectionTaskFeatureResponse struct {
}

type GetInspectionTaskFeatureResponse struct {
	Features []inspection.FeatureListItem `json:"features"`
}

type PostInspectionTaskDryRunRequest = map[string]any
