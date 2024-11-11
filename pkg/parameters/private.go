// Copyright 2024 Google LLC
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

package parameters

import (
	"strings"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/flag"
)

var Private *PrivateParameters = &PrivateParameters{}

type PrivateParameters struct {
	// InspectionMode
	// If this flag is set, KHI use the inspection token to access APIs.
	InspectionMode *bool

	// IAMToken is the initial IAM token used for accessing APIs.
	// InspectionMode will be true when only this value is specified.
	IAMToken *string

	// GALabels is the additional labels sent with the other labels. This flag is needed for taking analytics.
	// The expected format is "key1=value1,key2=value2,..."
	GALabels *string
}

// PostProcess implements ParameterStore.
func (p *PrivateParameters) PostProcess() error {
	if p.IAMToken != nil && *p.IAMToken != "" {
		*p.InspectionMode = true
	}
	return nil
}

// Prepare implements ParameterStore.
func (p *PrivateParameters) Prepare() error {
	p.InspectionMode = flag.Bool("inspection-mode", false, "If this flag is set, KHI use the inspection token to access APIs.", "")
	p.IAMToken = flag.String("iam-token", "", "The initial IAM token used for accessing APIs. InspectionMode will be true when only this value is specified.", "")
	p.GALabels = flag.String("ga-labels", "", "The additional labels sent with the other labels. This flag is needed for taking analytics.", "")
	return nil
}

// GetMapOfGALabels convert the GALabels parameter into a map.
// The expected format is "key1=value1,key2=value2,..."
func (p *PrivateParameters) GetMapOfGALabels() map[string]string {
	result := make(map[string]string)
	keyValuePairs := strings.Split(*p.GALabels, ",")
	for _, pair := range keyValuePairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		key, value, found := strings.Cut(pair, "=")
		if !found {
			value = "null"
		}
		result[key] = value
	}
	return result
}

var _ ParameterStore = (*PrivateParameters)(nil)
