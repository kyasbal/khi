// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package privatecomposer_contract

import (
	"math"

	coreinspection "github.com/GoogleCloudPlatform/khi/pkg/core/inspection"
	googlecloudcommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudcommon/contract"
	inspectioncore_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore/contract"
	privatecommon_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/privatecommon/contract"
)

// InspectionTypeId is the inspection type id for Private Cloud Composer v3.
var InspectionTypeId = "gcp-composer-v3"

var ComposerV3InspectionType = coreinspection.InspectionType{
	Id:          InspectionTypeId,
	Name:        "Managed Airflow (gen 3) (Internal)",
	Description: "Fetch both Managed Airflow (gen 3) backend and GKE tenant project logs to generate Managed Airflow (gen 3) specific results.",
	Icon:        "assets/icons/composer.webp",
	Priority:    math.MaxInt - 5,
	Labels: map[string]string{
		inspectioncore_contract.InspectionTypeLabelKeyLogSource:      "cloud_logging",
		inspectioncore_contract.InspectionTypeLabelKeyEnvironment:    "googlecloud",
		inspectioncore_contract.InspectionTypeLabelKeyBasePlatform:   "kubernetes",
		googlecloudcommon_contract.InspectionTypeLabelKeyClusterType: "gke",
		googlecloudcommon_contract.InspectionTypeLabelKeyProduct:     "composer",
		privatecommon_contract.InspectionTypeLabelKeyPrivate:         "true",
	},
}
