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

package index

import "github.com/GoogleCloudPlatform/khi/pkg/server/index"

// GoogleIdentityServiceTagGenerator returns its client tag. This is needed for Google Drive integration only for internal.
type GoogleIdentityServiceTagGenerator struct {
}

func (g *GoogleIdentityServiceTagGenerator) GenerateTags() []string {
	return []string{`<script src="https://accounts.google.com/gsi/client"></script>`}
}

var _ index.IndexTagGenerator = (*GoogleIdentityServiceTagGenerator)(nil)
