package aws

import (
	"math"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
)

var InspectionTypeId = "gcp-gke-on-aws"

var AnthosOnAWSInspectionType = inspection.InspectionType{
	Id:   InspectionTypeId,
	Name: "GKE on AWS(Anthos on AWS)",
	Description: `Visualize logs generated from GKE on AWS cluster. 
Supporting K8s audit log, k8s event log,k8s node log, k8s container log and MultiCloud API audit log.`,
	Icon:     "assets/icons/anthos.png",
	Priority: math.MaxInt - 2,
}
