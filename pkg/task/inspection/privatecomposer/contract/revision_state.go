// Copyright 2026 Google LLC
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
	pb "github.com/GoogleCloudPlatform/khi/pkg/generated/khifile/v6"
	"github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6/style"
)

var (
	// RevisionStateCloudSQLInstanceProvisioning represents the provisioning state of a Cloud SQL instance.
	RevisionStateCloudSQLInstanceProvisioning = style.MustRegisterRevisionState(
		"Cloud SQL Instance is being provisioned",
		"deployed_code_history",
		"The Cloud SQL database instance is currently being provisioned.",
		style.MustForceConvertSRGBHex("#6666ff"),
		pb.RevisionStateStyle_REVISION_STATE_STYLE_NORMAL,
	)
	// RevisionStateCloudSQLInstanceExisting represents the existing state of a Cloud SQL instance.
	RevisionStateCloudSQLInstanceExisting = style.MustRegisterRevisionState(
		"Cloud SQL Instance exists",
		"database",
		"The Cloud SQL database instance exists and is active.",
		style.Color{R: 0.0, G: 0.0, B: 1.0, A: 1.0},
		pb.RevisionStateStyle_REVISION_STATE_STYLE_NORMAL,
	)
	// RevisionStateCloudSQLInstanceDeleting represents the deleting state of a Cloud SQL instance.
	RevisionStateCloudSQLInstanceDeleting = style.MustRegisterRevisionState(
		"Cloud SQL Instance is being deleted",
		"auto_delete",
		"The Cloud SQL database instance is undergoing deletion.",
		style.Color{R: 0.8, G: 0.33333334, B: 0.0, A: 1.0},
		pb.RevisionStateStyle_REVISION_STATE_STYLE_NORMAL,
	)
	// RevisionStateCloudSQLInstanceDeleted represents the deleted state of a Cloud SQL instance.
	RevisionStateCloudSQLInstanceDeleted = style.MustRegisterRevisionState(
		"Cloud SQL Instance is deleted",
		"delete_forever",
		"The Cloud SQL database instance has been deleted.",
		style.Color{R: 0.8, G: 0.0, B: 0.0, A: 1.0},
		pb.RevisionStateStyle_REVISION_STATE_STYLE_DELETED,
	)
	// RevisionStateCloudSQLInstanceProvisioningLogNotFound represents provisioning with missing starting log.
	RevisionStateCloudSQLInstanceProvisioningLogNotFound = style.MustRegisterRevisionState(
		"Cloud SQL Instance is being provisioned, but starting log not found",
		"deployed_code_history",
		"The Cloud SQL database instance provisioning was started, but the starting log entry was not found in the selected time range.",
		style.MustForceConvertSRGBHex("#6666ff"),
		pb.RevisionStateStyle_REVISION_STATE_STYLE_PARTIAL_INFO,
	)
	// RevisionStateCloudSQLInstanceExistingLogNotFound represents existing instance with missing creation log.
	RevisionStateCloudSQLInstanceExistingLogNotFound = style.MustRegisterRevisionState(
		"Cloud SQL Instance exists, but creation log not found",
		"database",
		"The Cloud SQL database instance exists, but the creation or existence log entry was not found in the selected time range.",
		style.Color{R: 0.0, G: 0.0, B: 1.0, A: 1.0},
		pb.RevisionStateStyle_REVISION_STATE_STYLE_PARTIAL_INFO,
	)
	// RevisionStateCloudSQLInstanceDeletingLogNotFound represents deleting instance with missing starting log.
	RevisionStateCloudSQLInstanceDeletingLogNotFound = style.MustRegisterRevisionState(
		"Cloud SQL Instance is being deleted, but starting log not found",
		"auto_delete",
		"The Cloud SQL database instance deletion was in progress, but the deletion starting log entry was not found in the selected time range.",
		style.Color{R: 0.8, G: 0.33333334, B: 0.0, A: 1.0},
		pb.RevisionStateStyle_REVISION_STATE_STYLE_PARTIAL_INFO,
	)
)
