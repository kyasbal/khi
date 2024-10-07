package gke

import (
	"math"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
)

var InspectionTypeId = "gcp-gke"

var GKEInspectionType = inspection.InspectionType{
	Id:   InspectionTypeId,
	Name: "Google Kubernetes Engine",
	Description: `Visualize logs generated from GKE cluster. 
Supporting K8s audit log, K8s event log,K8s node log, K8s container log, GCE audit log, Networking audit log(NEG attach/detach) and autoscaler log.`,
	Icon:     "/assets/icons/gke.png",
	Priority: math.MaxInt,
}
