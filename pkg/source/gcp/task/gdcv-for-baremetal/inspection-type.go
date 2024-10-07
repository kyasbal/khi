package baremetal

import (
	"math"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
)

var InspectionTypeId = "gcp-gdcv-for-baremetal"

var AnthosOnBaremetalInspectionType = inspection.InspectionType{
	Id:   InspectionTypeId,
	Name: "GDCV for Baremetal(GKE on Baremetal, Anthos on Baremetal)",
	Description: `Visualize logs generated from GDCV for baremetal cluster(including user cluster/admin cluster/hybrid cluster or standalone cluster).
Supporting K8s audit log, k8s event log,k8s node log, k8s container log and OnPream API audit log.

This type can also be used for GCDE or GDCH.`,
	Icon:     "/assets/icons/anthos.png",
	Priority: math.MaxInt - 3,
}
