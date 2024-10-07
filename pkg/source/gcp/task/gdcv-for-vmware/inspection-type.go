package vmware

import (
	"math"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
)

var InspectionTypeId = "gcp-gdcv-for-vmware"

var AnthosOnVMWareInspectionType = inspection.InspectionType{
	Id:   InspectionTypeId,
	Name: "GDCV for VMWare(GKE on VMWare, Anthos on VMWare)",
	Description: `Visualize logs generated from GDCV for VMWare cluster(including admin clsuter/user cluster).
Supporting K8s audit log, k8s event log,k8s node log, k8s container log and OnPream API audit log.`,
	Icon:     "/assets/icons/anthos.png",
	Priority: math.MaxInt - 4,
}
