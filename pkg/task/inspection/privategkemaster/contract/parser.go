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

import (
	"github.com/GoogleCloudPlatform/khi/pkg/core/inspection/logutil"
	googlecloudlogk8scontrolplane_contract "github.com/GoogleCloudPlatform/khi/pkg/task/inspection/googlecloudlogk8scontrolplane/contract"
)

// DefaultPrivateGKEMasterControlplaneLogParser is the default SelectorLogParser for private GKE master control plane logs.
var DefaultPrivateGKEMasterControlplaneLogParser = logutil.NewSelectorLogParser[googlecloudlogk8scontrolplane_contract.ControlplaneLogContext](
	logutil.NewMultiTextLogParser(
		logutil.NewKLogTextParser(true),
		logutil.NewKLogTextParser(false),
		logutil.NewJsonlTextParser(),
		logutil.NewZapConsoleTextParser(),
		logutil.NewLogfmtTextParser(),
		&logutil.FallbackRawTextLogParser{},
	),
	logutil.ParserRule[googlecloudlogk8scontrolplane_contract.ControlplaneLogContext]{
		Match: func(ctx googlecloudlogk8scontrolplane_contract.ControlplaneLogContext) bool {
			return ctx.ComponentName == "scheduler" || ctx.ComponentName == "controller-manager"
		},
		Parser: logutil.NewMultiTextLogParser(
			logutil.NewKLogTextParser(true),
			logutil.NewKLogTextParser(false),
		),
	},
)
