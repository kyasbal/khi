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
	"errors"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common/flag"
)

// ParameterStore parses subset of parameters given to KHI program.
type ParameterStore interface {
	// Prepare initialize settings for reading parameters with `flag` package.
	Prepare() error

	// PostProcess override parsed parameters depending on the other parsed parameters.
	PostProcess() error
}

// Parse initializes the given parameter stores.
func Parse(stores ...ParameterStore) error {
	if flag.Parsed() {
		return errors.New("parameter flags are already parsed")
	}
	for _, store := range stores {
		err := store.Prepare()
		if err != nil {
			return err
		}
	}
	flag.Parse()
	for _, store := range stores {
		err := store.PostProcess()
		if err != nil {
			return err
		}
	}
	return nil
}
