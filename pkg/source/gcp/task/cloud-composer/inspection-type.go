package composer_task

import (
	"math"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
)

var InspectionTypeId = "gcp-composer"

var ComposerInspectionType = inspection.InspectionType{
	Id:   InspectionTypeId,
	Name: "Cloud Composer",
	Description: `Visualize logs related to Cloud Composer environment.
Supports all GKE related logs(Cloud Composer v2 or v1) and Airflow logs(Airflow 2.0.0 or higher in any Cloud Composer version(v1-v3))`,
	Icon:     "/assets/icons/composer.webp",
	Priority: math.MaxInt - 1,
}
