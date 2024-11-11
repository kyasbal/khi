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

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/flag"

var Common *CommonParameters = &CommonParameters{}

type CommonParameters struct {
	// DataDestinationFolder is the folder path where the final khi file to be stored for serving.
	DataDestinationFolder *string
	// TemporaryFolder is the folder path where be used as a working directory to generate the final khi file.
	TemporaryFolder *string
}

// PostProcess implements ParameterStore.
func (c *CommonParameters) PostProcess() error {
	return nil
}

// Prepare implements ParameterStore.
func (c *CommonParameters) Prepare() error {
	c.DataDestinationFolder = flag.String("data-destination-folder", "./data", "The folder path where the final khi file to be stored for serving.", "")
	c.TemporaryFolder = flag.String("temporary-folder", "/tmp", "The folder path where be used as a working directory to generate the final khi file.", "")
	return nil
}

var _ ParameterStore = (*CommonParameters)(nil)
