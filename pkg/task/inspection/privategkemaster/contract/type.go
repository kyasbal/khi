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

package privategkemaster_contract

// LogSource represents the source of the GKE Master logs.
type LogSource struct {
	// TenantProjectID is the project ID of the GKE Master project.
	TenantProjectID string
	// LogViewResourceName is the resource name of the log view.
	LogViewResourceName string
}
