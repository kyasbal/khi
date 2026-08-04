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
	"github.com/GoogleCloudPlatform/khi/pkg/model/khifile/v6/style"
)

var (
	// LogTypeCloudSQL is the style for Cloud SQL database engine logs in the tenant project.
	LogTypeCloudSQL = style.MustRegisterLogType(
		"composer-cloudsql",
		"Cloud Composer Tenant Cloud SQL Logs",
		style.MustForceConvertSRGBHex("#4285F4"),
		style.ColorWhite,
	)

	// LogTypeCloudSQLAudit is the style for Cloud SQL activity audit logs in the tenant project.
	LogTypeCloudSQLAudit = style.MustRegisterLogType(
		"composer-cloudsql-audit",
		"Cloud Composer Tenant Cloud SQL Audit Logs",
		style.MustForceConvertSRGBHex("#34A853"),
		style.ColorWhite,
	)
)
