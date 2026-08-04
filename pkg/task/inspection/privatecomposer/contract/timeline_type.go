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
	// TimelineTypeCloudSQLInstance is the style for an individual Cloud SQL database instance.
	TimelineTypeCloudSQLInstance = style.MustRegisterTimelineType(
		"Cloud SQL",
		"Cloud SQL Instance",
		"database",
		1,
		style.Color{R: 0.780, G: 0.863, B: 1.000, A: 1.0},
		style.ColorBlack,
		style.Color{R: 0.259, G: 0.522, B: 0.957, A: 1.0},
		style.ColorWhite,
		true,
		4001,
		style.AlphabeticalSortPolicy(),
	)

	// TimelineTypeCloudSQLLog is the style for a Cloud SQL database engine log file.
	TimelineTypeCloudSQLLog = style.MustRegisterTimelineType(
		"Cloud SQL logs",
		"Cloud SQL Log",
		"notes",
		0.6,
		style.ColorWhite,
		style.ColorBlack,
		style.MustForceConvertSRGBHex("#34A853"),
		style.ColorWhite,
		true,
		4010,
		style.AlphabeticalSortPolicy(),
	)
)
